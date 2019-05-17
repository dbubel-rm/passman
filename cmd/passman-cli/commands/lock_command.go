package commands

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

	password := CleanInput(string(bytePassword))
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
	return "Ex) passman lock -  use an easy to remember password that is still strong"
}

func (c *LockCommand) Synopsis() string {
	return "lock the config file"
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func CleanInput(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
