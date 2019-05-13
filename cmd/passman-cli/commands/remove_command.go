package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/cli"
)

type RemoveCommand struct {
	ServiceName string
	UI          cli.Ui
}

func (c *RemoveCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("remove", flag.ContinueOnError)
	cmdFlags.StringVar(&c.ServiceName, "s", "", "Service name")
	cmdFlags.Parse(args)

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
	req, err := http.NewRequest("DELETE", url, nil)
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

	c.UI.Output(fmt.Sprintf("Credentials for service %s deleted OK", c.ServiceName))
	return 0
}

func (c *RemoveCommand) Help() string {
	return "passman remove -s <service_name>"
}

func (c *RemoveCommand) Synopsis() string {
	return "Remove a credential from Passman."
}
