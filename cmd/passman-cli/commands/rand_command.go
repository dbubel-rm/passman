package commands

import (
	"flag"
	"os"

	"github.com/dbubel/passman/cmd/passman-cli-legacy/utils"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
)

type RandCommand struct {
	Length int
	UI     cli.Ui
}

func (c *RandCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("rand", flag.ContinueOnError)
	cmdFlags.IntVar(&c.Length, "l", 0, "Length")
	cmdFlags.Parse(args)

	if c.Length == 0 {
		c.UI.Warn(c.Help())
		return 1
	}

	n, _ := utils.GenerateRandomString(c.Length)

	data := [][]string{
		{n},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"New Password"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	return 0
}

func (c *RandCommand) Help() string {
	return "passman rand -l <length>"
}

func (c *RandCommand) Synopsis() string {
	return "Generate a random password"
}
