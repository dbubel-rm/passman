package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/dbubel/passman/cmd/passman-cli/models"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
)

type ListCommand struct {
	UI cli.Ui
}

func (c *ListCommand) Run(args []string) int {

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

	url := cfg.Backend + "/v1/services"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.IDToken))
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

	var credentialRecord []models.Credential

	err = json.Unmarshal(body, &credentialRecord)

	fmt.Println(string(body))

	data := [][]string{}
	for i := range credentialRecord {

		t, err := time.Parse("2006-01-02 15:04:05", credentialRecord[i].UpdatedAt)
		s := fmt.Sprintf("%v", -1*math.Round(t.Sub(time.Now()).Hours()/24))
		credentialRecord[i].UpdatedAt = s

		if err != nil {
			log.Println(err.Error())
		}

		data = append(data, []string{credentialRecord[i].CredentialID, credentialRecord[i].ServiceName, credentialRecord[i].UpdatedAt})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Service", "Last Updated day(s)"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	return 0
}

func (c *ListCommand) Help() string {
	return "Ex) passman list"
}

func (c *ListCommand) Synopsis() string {
	return "List all the services you've stored in passman"
}

// type credentialRecords []struct {
// 	CredentialID string
// 	ServiceName  string
// 	UpdatedAt    string
// }
