package commands

import (
	"encoding/json"
	"fmt"
	"github.com/dbubel/passman/cmd/passman-cli/models"
	"github.com/dbubel/passman/cmd/passman-cli/utils"
	"io/ioutil"
	"log"
	"net/http"
)

func Rm(argsWithoutProg []string) {
	if len(argsWithoutProg) != 2 {
		log.Println("Not enough arguments")
		return
	}

	tokenData, err := utils.GetUserStore(PassmanHome)

	if err != nil {
		log.Println(err.Error())
		return
	}

	var storedJWT models.FirebaseStruct
	err = json.Unmarshal(tokenData, &storedJWT)

	if err != nil {
		log.Println(err.Error())
	}

	serviceName := argsWithoutProg[1]
	req, err := http.NewRequest("DELETE", urlNewCredential+"/"+serviceName, nil)
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
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println(string(body))
		return
	}
	log.Println("Credential deleted OK:")
}
