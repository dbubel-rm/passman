package commands

import (
	"encoding/json"
	"fmt"
	"github.com/dbubel/passman/cmd/passman-cli/models"
	"github.com/dbubel/passman/cmd/passman-cli/secret"
	"github.com/dbubel/passman/cmd/passman-cli/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func Insert(argsWithoutProg []string) {
	if len(argsWithoutProg) != 4 {
		log.Println("Not enough args wanted 4 got", len(argsWithoutProg))
		return
	}

	newCredentialPayload := `{"serviceName": "%s","username": "%s","password": "%s"}`
	serviceName := ""
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

	username, password := argsWithoutProg[2], argsWithoutProg[3]

	password = secret.Encrypt([]byte(password), os.Getenv(PASSMAN_MASTER))
	username = secret.Encrypt([]byte(username), os.Getenv(PASSMAN_MASTER))
	serviceName = argsWithoutProg[1]
	newCredentialPayload = fmt.Sprintf(newCredentialPayload, serviceName, username, password)
	req, err := http.NewRequest("POST", urlNewCredential, strings.NewReader(newCredentialPayload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storedJWT.IDToken))

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

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println(string(body))
		log.Println("Try logging in again")
		return
	}

	log.Println("Credential added OK")
}
