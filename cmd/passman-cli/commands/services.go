package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

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

		t, err := time.Parse("2006-01-02 15:04:05", credentialRecord[i].UpdatedAt)
		s := fmt.Sprintf("%v", -1*math.Round(t.Sub(time.Now()).Hours()/24))
		credentialRecord[i].UpdatedAt = s

		if err != nil {
			fmt.Println(err.Error())
		}
		// fmt.Println(t)
		data = append(data, []string{credentialRecord[i].CredentialID, credentialRecord[i].ServiceName, credentialRecord[i].UpdatedAt})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Service", "Last Updated day(s)"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
