package commands

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/cli"
	"golang.org/x/crypto/ssh/terminal"
)

type UnlockCommand struct {
	Password string
	UserHome string
	UI       cli.Ui
}

func (c *UnlockCommand) Run(args []string) int {

	dat, err := ioutil.ReadFile(c.UserHome + "/.passman/config.json")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	password := CleanInput(string(bytePassword))
	fmt.Println("")

	sDec, err := b64.StdEncoding.DecodeString(string(dat))

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	dec, err := Decrypt(sDec, password)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	err = ioutil.WriteFile(c.UserHome+"/.passman/config.json", []byte(dec), 0644)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output("Decrypted OK")
	return 0
}

func (c *UnlockCommand) Help() string {
	return "Ex) passman unlock"
}

func (c *UnlockCommand) Synopsis() string {
	return "Unlock the config file"
}
