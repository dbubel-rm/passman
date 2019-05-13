package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/dbubel/passman/cmd/passman-cli/commands"
	"github.com/dbubel/passman/cmd/passman-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

//var baseUrl = "https://ec2-100-25-42-237.compute-1.amazonaws.com:3000"

var argsWithoutProg = os.Args[1:]

func version(a []string) {
	fmt.Println("v0.0.5")
}

func main() {

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	commands.PassmanHome = usr.HomeDir + "/.passman/session.json"

	if os.Getenv(commands.PASSMAN_MASTER) == "" {
		master := getUsernameAndPassword()
		os.Setenv(commands.PASSMAN_MASTER, master)
	}

	log.SetFlags(log.Lshortfile)

	actions := make(map[string]func(a []string))
	actions[commands.HELP] = commands.Help
	actions[commands.VERSION] = version
	actions[commands.GEN_PASS] = commands.GenPassword
	// API calls
	actions[commands.LOGIN] = commands.Signin
	actions[commands.REGISTER_ACCOUNT] = commands.Register
	actions[commands.NUKE_ACCOUNT] = commands.Nuke
	actions[commands.INSERT_CREDENTIAL] = commands.Insert
	actions[commands.GET_CREDENTIAL] = commands.Get
	actions[commands.RM_CREDENTIAL] = commands.Rm
	actions[commands.GET_SERVICES] = commands.Services
	actions[commands.UPDATE_CREDENTIAL] = commands.Update
	actions[commands.UPDATE_MASTER] = commands.UpdateMasterPass

	if len(argsWithoutProg) == 0 {
		log.Println("No action specified")
		commands.Help(argsWithoutProg)
		return
	}

	_, err = os.Stat(commands.PassmanHome)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	if os.IsNotExist(err) {
		fmt.Println("trying to create")
		_, e := os.Create(commands.PassmanHome)
		if e != nil {
			fmt.Println("bad create", e.Error())
		}
	}

	action, ok := actions[argsWithoutProg[0]]
	if ok {
		action(argsWithoutProg)
	} else {
		log.Println("Invalid action specified")
		commands.Help(argsWithoutProg)
	}
}

func getUsernameAndPassword() string {
	// fmt.Print("Username: ")
	// text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	// text = cleanInput(text)
	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	password := utils.CleanInput(string(bytePassword))
	fmt.Println("")
	return password
}
