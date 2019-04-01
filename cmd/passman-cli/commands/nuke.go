package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/dbubel/passman/cmd/passman-cli/models"
	"github.com/dbubel/passman/cmd/passman-cli/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var urlDeleteUser = baseUrl + "/v1/users"

func Nuke(argsWithoutProg []string) {
	tokenData, err := utils.GetUserStore(PassmanHome)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var storedJWT models.FirebaseStruct
	err = json.Unmarshal(tokenData, &storedJWT)

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("!!! WARNING !!!")
	log.Println("Are you sure you want to delete ALL data associaed with:", storedJWT.Email)
	log.Println("Type: 'yes' to remove all data")

	text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	text = utils.CleanInput(text)
	if text != "yes" {
		log.Println("Aborting")
		return
	}

	var payload = `{"idToken":"%s"}`
	payload = fmt.Sprintf(payload, storedJWT.IDToken)

	req, err := http.NewRequest("DELETE", urlDeleteUser, strings.NewReader(payload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storedJWT.IDToken))

	if err != nil {
		log.Println(err.Error())
	}

	res, err := utils.SkipTLS(req)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	log.Println(string(body))
	os.Remove(PassmanHome)
}
