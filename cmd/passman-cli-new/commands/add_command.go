package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dbubel/passman/cmd/passman-cli/secret"
	"github.com/mitchellh/cli"
)

type AddCommand struct {
	ServiceName string
	UserName    string
	Password    string
	UI          cli.Ui
}

func (c *AddCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("add", flag.ContinueOnError)
	cmdFlags.StringVar(&c.ServiceName, "service", "", "Service name")
	cmdFlags.StringVar(&c.UserName, "p", "", "Username")
	cmdFlags.StringVar(&c.Password, "u", "", "Password")
	cmdFlags.Parse(args)

	if c.ServiceName == "" || c.Password == "" || c.UserName == "" {
		c.UI.Error(c.Help())
		return 1
	}

	cfg, err := getConfigInfo()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	session, err := getSessionInfo()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	newCredentialPayload := `{"serviceName": "%s","username": "%s","password": "%s"}`

	c.Password = secret.Encrypt([]byte(c.Password), cfg.Password)
	c.UserName = secret.Encrypt([]byte(c.UserName), cfg.Password)
	// c.ServiceName = secret.Encrypt([]byte(c.ServiceName), cfg.Password)

	newCredentialPayload = fmt.Sprintf(newCredentialPayload, c.ServiceName, c.UserName, c.Password)
	req, err := http.NewRequest("POST", cfg.Backend+"/v1/credential", strings.NewReader(newCredentialPayload))
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

	c.UI.Output("Credential inserted OK")
	return 0
}

func (c *AddCommand) Help() string {
	return "Ex) passman -u <username> -p <password> -s <service>"
}

func (c *AddCommand) Synopsis() string {
	return "Adds a credential to you passman datastore"
}
