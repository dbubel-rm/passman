package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"strings"

	"github.com/dbubel/passman/cmd/passman-cli-new/models"
	"github.com/mitchellh/cli"
)

func getConfigInfo() (models.Config, error) {
	var cfg models.Config
	usr, err := user.Current()
	if err != nil {
		return cfg, err
	}

	passmanConfig := usr.HomeDir + "/.passman/config.json"
	dat, err := ioutil.ReadFile(passmanConfig)
	json.Unmarshal(dat, &cfg)
	return cfg, err
}

type LoginCommand struct {
	EmailAddress string
	Ui           cli.Ui
	UserHome     string
}

func (c *LoginCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("login", flag.ContinueOnError)
	cmdFlags.StringVar(&c.EmailAddress, "u", "", "Email address")
	cmdFlags.Parse(args)

	if c.EmailAddress == "" {
		c.Ui.Error(fmt.Sprint(c.Help()))
		return 1
	}

	cfg, err := getConfigInfo()

	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	var payload = `{"email":"%s","password":"%s","returnSecureToken": true}`
	payload = fmt.Sprintf(payload, c.EmailAddress, cfg.Password)
	req, _ := http.NewRequest("GET", cfg.Backend+"/v1/signin", strings.NewReader(payload))

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Println(err.Error())
		return 1
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if res.StatusCode != 200 {
		c.Ui.Error(string(body))
		return 1
	}

	err = ioutil.WriteFile(c.UserHome+"/.passman/session.json", body, 0644)

	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprint("Login OK"))
	return 0
}

func (c *LoginCommand) Help() string {
	return "Login to your passman account (Ex. passman login -u <example@email.com>)"
}

func (c *LoginCommand) Synopsis() string {
	return "Login to passman"
}
