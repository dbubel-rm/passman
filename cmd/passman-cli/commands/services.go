package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dbubel/passman/cmd/passman-cli/models"
	"github.com/dbubel/passman/cmd/passman-cli/utils"
	"github.com/olekukonko/tablewriter"
)

var urlServices = baseURL + "/v1/services"

func Services(argsWithoutProgs []string) {
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

	req, err := http.NewRequest("GET", urlServices+"/", nil)
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

	var credentialRecord = []struct {
		CredentialID string
		ServiceName  string
		UpdatedAt    string
	}{}

	err = json.Unmarshal(body, &credentialRecord)

	if err != nil {
		log.Println(err.Error())
		return
	}

	data := [][]string{}
	for i := range credentialRecord {
		data = append(data, []string{credentialRecord[i].ServiceName, credentialRecord[i].CredentialID, credentialRecord[i].UpdatedAt})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Service", "ID", "Updated At"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
