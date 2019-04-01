package commands

import (
	"encoding/json"
	"fmt"
	"github.com/dbubel/passman/cmd/passman-cli/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var urlCreateAccount = baseUrl + "/v1/users"
var urlVerifyAccount = baseUrl + "/v1/users/verify"

func Register(argsWithoutProg []string) {
	payloadCreateAccount := `{"email": "%s","password": "%s","returnSecureToken": true}`
	if len(argsWithoutProg) != 2 {
		fmt.Println("Invalid option")
		Help(argsWithoutProg)
		return
	}
	// fmt.Println(argsWithoutProg)
	username := argsWithoutProg[1]
	password := os.Getenv(PASSMAN_MASTER)

	payloadCreateAccount = fmt.Sprintf(payloadCreateAccount, username, password)
	req, err := http.NewRequest("POST", urlCreateAccount, strings.NewReader(payloadCreateAccount))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res, err := utils.SkipTLS(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println(string(body))
		return
	}

	err = ioutil.WriteFile(PassmanHome, body, 0644)

	if err != nil {
		log.Println(err.Error())
		return
	}
	type response struct {
		IdToken string `json:"idToken"`
	}
	var id response
	json.Unmarshal(body, &id)

	payloadVerifyAccount := fmt.Sprintf(`{"requestType": "VERIFY_EMAIL","idToken": "%s"}`, id.IdToken)
	req, err = http.NewRequest("POST", urlVerifyAccount, strings.NewReader(payloadVerifyAccount))
	res, err = utils.SkipTLS(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	body, _ = ioutil.ReadAll(res.Body)

	log.Println("Account created OK, check you email for a verification link.")
}
