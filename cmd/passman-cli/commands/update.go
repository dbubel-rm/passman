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

var urlUpdateCredential = baseURL + "/v1/credential/update"

func Update(argsWithoutProg []string) {
	if len(argsWithoutProg) != 3 {
		log.Println("Not enough args")
		return
	}

	newCredentialPayload := `{"serviceName": "%s","password": "%s"}`

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

	username, password := argsWithoutProg[1], argsWithoutProg[2]
	password = secret.Encrypt([]byte(password), os.Getenv(PASSMAN_MASTER))
	newCredentialPayload = fmt.Sprintf(newCredentialPayload, username, password)
	req, err := http.NewRequest("POST", urlUpdateCredential, strings.NewReader(newCredentialPayload))
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

	body, _ := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println(string(body))
		return
	}

	log.Println("Credential updated OK")
}
