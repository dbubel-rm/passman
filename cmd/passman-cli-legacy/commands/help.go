package commands

import "fmt"

func Help(a []string) {
	fmt.Println("Passman is a utility for managing your passwords.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\tpassman <command> [arguments]")
	fmt.Println("")
	fmt.Println("The commands are:")
	fmt.Println("")
	fmt.Printf("\t%s\t\tCreates a new passman account. Ex) passman %s newexample@example.com\n", REGISTER_ACCOUNT, REGISTER_ACCOUNT)
	fmt.Printf("\t%s\t\tInsert a credential. Ex) passman %s serviceName username password\n", INSERT_CREDENTIAL, INSERT_CREDENTIAL)
	fmt.Printf("\t%s\t\tGet a stored credential. Ex) passman %s serviceName\n", GET_CREDENTIAL, GET_CREDENTIAL)
	fmt.Printf("\t%s\t\tDeletes a stored credential. Ex) passman %s service_name\n", RM_CREDENTIAL, RM_CREDENTIAL)
	fmt.Printf("\t%s\t\tDeletes ALL credentials saved under you active account. Ex) passman %s\n", NUKE_ACCOUNT, NUKE_ACCOUNT)
	fmt.Printf("\t%s\t\tAuthenticate a passman session good for 30 minutes\n", LOGIN)
	fmt.Printf("\t%s\t\tGenerates a crypto random string. Ex) passman rand 16\n", GEN_PASS)
	fmt.Printf("\t%s\t\tGets all the services in the accout\n", GET_SERVICES)
	fmt.Printf("\t%s\t\tDisplays this message\n", HELP)
	fmt.Printf("\t%s\t\tDisplays the version of passman\n", VERSION)
}
