package commands

import (
	"fmt"
	"github.com/dbubel/passman/cmd/passman-cli/utils"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"strconv"
)

func GenPassword(argsWithoutProg []string) {
	if len(argsWithoutProg) != 2 {
		log.Println("length not specified")
		return
	}
	i, err := strconv.ParseInt(argsWithoutProg[1], 10, 32)
	if err != nil {
		fmt.Println(err.Error())
	}
	n, _ := utils.GenerateRandomString(int(i))

	data := [][]string{
		[]string{n},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"New Password"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
