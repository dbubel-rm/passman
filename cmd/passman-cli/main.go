package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dbubel/passman/cmd/passman-cli/commands"
	"github.com/mitchellh/cli"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func main() {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := cli.NewCLI("passman cli", "1.0.0")

	usr, _ := user.Current()
	c.Args = os.Args[1:]
	// c.Autocomplete = true

	c.Commands = map[string]cli.CommandFactory{
		"login": func() (cli.Command, error) {
			return &commands.LoginCommand{
				Ui: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorBlue,
					ErrorColor:  cli.UiColorRed,
				},
				UserHome: usr.HomeDir,
			}, nil
		},
		"new": func() (cli.Command, error) {
			return &commands.SignupCommand{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorBlue,
					ErrorColor:  cli.UiColorRed,
				},
				UserHome: usr.HomeDir,
			}, nil
		},
		"lock": func() (cli.Command, error) {
			return &commands.LockCommand{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorBlue,
					ErrorColor:  cli.UiColorRed,
				},
				UserHome: usr.HomeDir,
			}, nil
		},
		"unlock": func() (cli.Command, error) {
			return &commands.UnlockCommand{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorBlue,
					ErrorColor:  cli.UiColorRed,
				},
				UserHome: usr.HomeDir,
			}, nil
		},
		"get": func() (cli.Command, error) {
			return &commands.GetCommand{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorBlue,
					ErrorColor:  cli.UiColorRed,
				},
			}, nil
		},
		"add": func() (cli.Command, error) {
			return &commands.AddCommand{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorBlue,
					ErrorColor:  cli.UiColorRed,
				},
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &commands.ListCommand{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorBlue,
					ErrorColor:  cli.UiColorRed,
				},
			}, nil
		},
		"remove": func() (cli.Command, error) {
			return &commands.RemoveCommand{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorBlue,
					ErrorColor:  cli.UiColorRed,
				},
			}, nil
		},
		"rand": func() (cli.Command, error) {
			return &commands.RandCommand{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorBlue,
					ErrorColor:  cli.UiColorRed,
					WarnColor:   cli.UiColorGreen,
				},
			}, nil
		},
	}

	_, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	// os.Exit(exitStatus)
}
