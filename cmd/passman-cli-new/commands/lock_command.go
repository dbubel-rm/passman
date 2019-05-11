package commands

import (
	b64 "encoding/base64"
	"fmt"

	"github.com/mitchellh/cli"
)

type LockCommand struct {
	Password string
	UserHome string
	UI       cli.Ui
}

func (c *LockCommand) Run(args []string) int {

	// Grab info from CLI
	// cmdFlags := flag.NewFlagSet("signup", flag.ContinueOnError)
	// cmdFlags.StringVar(&c.Username, "u", "", "Email address for new account")
	// cmdFlags.StringVar(&c.Password, "p", "", "Password for new account")
	// cmdFlags.StringVar(&c.Hostname, "hostname", "", "Passman API url")
	// cmdFlags.Parse(args)

	// if c.Username == "" {
	// 	c.UI.Error(fmt.Sprint(c.Help()))
	// 	return 1
	// }
	// if c.Password == "" {
	// 	c.UI.Error(fmt.Sprint(c.Help()))
	// 	return 1
	// }
	// if c.Hostname == "" {
	// 	c.UI.Error(fmt.Sprint(c.Help()))
	// 	return 1
	// }

	// // Write the config file
	// var w models.Config
	// w.Backend = c.Hostname
	// w.Password = c.Password
	// w.Username = c.Username

	// dat, _ := json.Marshal(w)
	// err := ioutil.WriteFile(c.UserHome+"/.passman/config.json", dat, 0644)

	// if err != nil {
	// 	c.UI.Error(err.Error())
	// 	return 1
	// }

	enc := Encrypt([]byte("testingtestingtestingtestingtestingtestingtestingtestingtesting 123"), "password123")
	fmt.Println(enc)
	sDec, _ := b64.StdEncoding.DecodeString(enc)
	dec := Decrypt([]byte(sDec), "password123")
	fmt.Println(dec)
	return 0

	// cfg, err := getConfigInfo()

	// if err != nil {
	// 	c.UI.Error(err.Error())
	// 	return 1
	// }

	// dat, err := json.Marshal(cfg)
	// fmt.Print("Password: ")
	// bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	// password := utils.CleanInput(string(bytePassword))
	// fmt.Println("")

	// enc := Encrypt(dat, password)

	// ioutil.WriteFile(c.UserHome+"/.passman/config.json", []byte(enc), 0644)

	// dat, err = ioutil.ReadFile(c.UserHome + "/.passman/config.json")
	// if err != nil {
	// 	c.UI.Error(err.Error())
	// 	return 1
	// }

	// dec := Decrypt(dat, password)
	// fmt.Println(string(dec))
	// ioutil.WriteFile(c.UserHome+"/.passman/config.json", []byte(dec), 0644)

	// if err != nil {
	// 	c.UI.Error(err.Error())
	// 	return 1
	// }

	// // Create the account
	// createAccountJSON := `{"email": "%s","password": "%s","returnSecureToken": true}`
	// createAccountJSON = fmt.Sprintf(createAccountJSON, c.Username, cfg.Password)
	// req, err := http.NewRequest("POST", cfg.Backend+"/v1/users", strings.NewReader(createAccountJSON))

	// httpClient := http.Client{}
	// res, err := httpClient.Do(req)

	// if err != nil {
	// 	log.Println(err.Error())
	// 	return 1
	// }

	// body, err := ioutil.ReadAll(res.Body)

	// if err != nil {
	// 	c.UI.Error(err.Error())
	// 	return 1
	// }

	// defer res.Body.Close()

	// if res.StatusCode != 200 {
	// 	c.UI.Error(string(body))
	// 	return 1
	// }

	// Store the session data
	// err = ioutil.WriteFile(c.UserHome+"/.passman/session.json", body, 0644)

	// if err != nil {
	// 	c.UI.Error(err.Error())
	// 	return 1
	// }

	// var id models.FirebaseSession
	// err = json.Unmarshal(body, &id)

	// if err != nil {
	// 	c.UI.Error(err.Error())
	// 	return 1
	// }

	// verifyAccountJSON := fmt.Sprintf(`{"requestType": "VERIFY_EMAIL","idToken": "%s"}`, id.IDToken)
	// req, err = http.NewRequest("POST", cfg.Backend+"/v1/users/verify", strings.NewReader(verifyAccountJSON))
	// res, err = httpClient.Do(req)

	// if err != nil {
	// 	c.UI.Error(err.Error())
	// 	return 1
	// }

	c.UI.Output("Encrypted OK")
	return 0
}

func (c *LockCommand) Help() string {
	return "Ex) passman lock"
}

func (c *LockCommand) Synopsis() string {
	return "lock the config file"
}
