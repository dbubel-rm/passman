package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dbubel/passman/cmd/passman-cli/utils"
	"github.com/mitchellh/cli"
	"golang.org/x/crypto/ssh/terminal"
)

type LockCommand struct {
	Password string
	UserHome string
	UI       cli.Ui
}

func (c *LockCommand) Run(args []string) int {
	cfg, err := getConfigInfo()

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	dat, err := json.Marshal(cfg)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	password := utils.CleanInput(string(bytePassword))
	fmt.Println("")

	enc := Encrypt(dat, password)

	err = ioutil.WriteFile(c.UserHome+"/.passman/config.json", []byte(enc), 0644)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output("Encrypted OK")
	return 0
}

func (c *LockCommand) Help() string {
	return "Ex) passman lock"
}

func (c *LockCommand) Synopsis() string {
	return "lock the config file"
}
