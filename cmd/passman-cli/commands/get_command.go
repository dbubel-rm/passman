package commands

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dbubel/passman/cmd/passman-cli/models"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
)

type GetCommand struct {
	ServiceName string
	UI          cli.Ui
}

func (c *GetCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("get", flag.ContinueOnError)
	cmdFlags.StringVar(&c.ServiceName, "s", "", "Service name")
	cmdFlags.Parse(args)

	if c.ServiceName == "" {
		c.UI.Error(c.Help())
		return 1
	}

	session, err := getSessionInfo()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	cfg, err := getConfigInfo()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	url := cfg.Backend + "/v1/credential/" + c.ServiceName
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.IDToken))
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if res.StatusCode != 200 {
		c.UI.Error(string(body))
		return 1
	}

	credentialRecord := models.Credential{}
	err = json.Unmarshal(body, &credentialRecord)
	fmt.Println(credentialRecord)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	decodedPassword, err := b64.StdEncoding.DecodeString(credentialRecord.Password)
	decodedUsername, _ := b64.StdEncoding.DecodeString(credentialRecord.Username)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Drcypt password
	decryptedPassword, err := Decrypt([]byte(decodedPassword), cfg.Password)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	credentialRecord.Password = string(decryptedPassword)

	// Decrypt username
	decryptedUsername, err := Decrypt([]byte(decodedUsername), cfg.Password)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	credentialRecord.Username = string(decryptedUsername)

	data := [][]string{
		[]string{c.ServiceName, credentialRecord.Username, credentialRecord.Password},
	}
	fmt.Println(c.ServiceName, credentialRecord.Username, credentialRecord.Password)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"service name", "username", "password"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	return 0
}

func (c *GetCommand) Help() string {
	return "passman get -s <service_name>"
}

func (c *GetCommand) Synopsis() string {
	return "Get a credential that you've stored in passman"
}
