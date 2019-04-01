package commands

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dbubel/passman/cmd/passman-cli/models"
	"github.com/dbubel/passman/cmd/passman-cli/secret"
	"github.com/dbubel/passman/cmd/passman-cli/utils"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var urlNewCredential = baseUrl + "/v1/credential"

func Get(argsWithoutProg []string) {
	if len(argsWithoutProg) != 2 {
		log.Println("No service name")
		return
	}

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
	}

	serviceName = argsWithoutProg[1]
	req, err := http.NewRequest("GET", urlNewCredential+"/"+serviceName, nil)
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

	var credentialRecord = struct {
		ServiceName string
		Username    string
		Password    string
	}{}

	err = json.Unmarshal(body, &credentialRecord)

	if err != nil {
		log.Println(err.Error())
		return
	}

	sDec, err := b64.StdEncoding.DecodeString(credentialRecord.Password)
	uName, _ := b64.StdEncoding.DecodeString(credentialRecord.Username)

	if err != nil {
		log.Println(err.Error())
		return
	}

	credentialRecord.Password = string(secret.Decrypt([]byte(sDec), os.Getenv(PASSMAN_MASTER)))
	credentialRecord.Username = string(secret.Decrypt([]byte(uName), os.Getenv(PASSMAN_MASTER)))

	data := [][]string{
		[]string{credentialRecord.ServiceName, credentialRecord.Username, credentialRecord.Password},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"service name", "username", "password"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}
