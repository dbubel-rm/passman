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

	"github.com/dbubel/passman/cmd/passman-cli/models"
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
	err = json.Unmarshal(dat, &cfg)
	return cfg, err
}

func getSessionInfo() (models.FirebaseSession, error) {
	var cfg models.FirebaseSession
	usr, err := user.Current()
	if err != nil {
		return cfg, err
	}

	session := usr.HomeDir + "/.passman/session.json"
	dat, err := ioutil.ReadFile(session)
	err = json.Unmarshal(dat, &cfg)
	return cfg, err
}

type LoginCommand struct {
	Username string
	Password string
	Hostname string
	Ui       cli.Ui
	UserHome string
}

func (c *LoginCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("login", flag.ContinueOnError)
	cmdFlags.StringVar(&c.Username, "u", "", "Email address")
	cmdFlags.StringVar(&c.Password, "p", "", "Password")
	cmdFlags.StringVar(&c.Hostname, "hostname", "", "Hostname")
	cmdFlags.Parse(args)

	cfg, _ := getConfigInfo()

	if c.Username == "" {
		c.Username = cfg.Username
	}

	if c.Username == "" {
		c.Ui.Error(c.Help())
		return 1
	}

	if c.Password == "" {
		c.Password = cfg.Password
	}

	if c.Password == "" {
		c.Ui.Error(c.Help())
		return 1
	}

	if c.Hostname == "" {
		c.Hostname = cfg.Backend
	}

	if c.Hostname == "" {
		c.Ui.Error(c.Help())
		return 1
	}

	payload := `{"email":"%s","password":"%s","returnSecureToken": true}`
	payload = fmt.Sprintf(payload, c.Username, c.Password)
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
	return "passman login -u <example@email.com> -p <password> -hostname <https://somehost.som>"
}

func (c *LoginCommand) Synopsis() string {
	return "Login to passman"
}
