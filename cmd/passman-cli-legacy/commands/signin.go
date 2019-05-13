package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/dbubel/passman/cmd/passman-cli/utils"
)

var urlAuthUser = baseURL + "/v1/signin"
var Usr string

func init() {
	u, _ := user.Current()
	Usr = u.HomeDir
}

func Signin(argsWithoutProg []string) {
	if len(argsWithoutProg) != 2 {
		log.Println("No account specified")
		return
	}

	username := argsWithoutProg[1]
	password := os.Getenv(PASSMAN_MASTER)

	fmt.Println("Using pass", password)
	var payload = `{"email":"%s","password":"%s","returnSecureToken": true}`
	payload = fmt.Sprintf(payload, username, password)
	req, err := http.NewRequest("GET", urlAuthUser, strings.NewReader(payload))

	if err != nil {
		log.Println(err.Error())
		return
	}

	res, err := utils.SkipTLS(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Println(err.Error())
	}

	if res.StatusCode != 200 {
		log.Println(string(body))
		return
	}

	log.Println("Login OK")

	err = ioutil.WriteFile(Usr+"/.passman/session.json", body, 0644)

	if err != nil {
		log.Println(err.Error())
	}
}
