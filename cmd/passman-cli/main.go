package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/bcrypt"
)

type FirebaseStruct struct {
	Kind           string `json:"kind"`
	LocalID        string `json:"localId"`
	Email          string `json:"email"`
	DisplayName    string `json:"displayName"`
	IDToken        string `json:"idToken"`
	Registered     bool   `json:"registered"`
	RefreshToken   string `json:"refreshToken"`
	ExpiresIn      string `json:"expiresIn"`
	PassmanPayload string `json:"passmanPayload"`
}

var baseUrl = "http://localhost:3000"

var urlCreateAccount = baseUrl + "/v1/users"
var urlDeleteUser = baseUrl + "/v1/users"
var urlAuthUser = baseUrl + "/v1/signin"
var urlNewCredential = baseUrl + "/v1/credential"
var urlUpdateCredential = baseUrl + "/v1/credential/update"
var urlServices = baseUrl + "/v1/services"
var urlVerifyAccount = baseUrl + "/v1/users/verify"
var passmanHome = "~/.passman/session.json"

var argsWithoutProg = os.Args[1:]

const (
	REGISTER_ACCOUNT  = "init"
	NUKE_ACCOUNT      = "nuke"
	HELP              = "help"
	VERSION           = "version"
	GEN_PASS          = "rand"
	INSERT_CREDENTIAL = "insert"
	LOGIN             = "login"
	PASSMAN_MASTER    = "PASSMAN_MASTER"
	GET_CREDENTIAL    = "get"
	// GET_CREDS      = "get:credentials"
	RM_CREDENTIAL     = "rm"
	GET_SERVICES      = "services"
	UPDATE_CREDENTIAL = "update"
)

func version() {
	fmt.Println("v0.0.1")

}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func help() {
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
	fmt.Printf("\t%s\t\tDisplays this message\n", HELP)
	fmt.Printf("\t%s\t\tDisplays the version of passman\n", VERSION)
}

func main() {

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	passmanHome = usr.HomeDir + "/.passman/session.json"

	if os.Getenv(PASSMAN_MASTER) == "" {
		log.Printf("No %s environment vairable set\n", PASSMAN_MASTER)
		return
	}

	log.SetFlags(log.Lshortfile)

	actions := make(map[string]func())
	actions[HELP] = help
	actions[VERSION] = version
	actions[GEN_PASS] = genPassword
	// API calls
	actions[LOGIN] = signin
	actions[REGISTER_ACCOUNT] = register
	actions[NUKE_ACCOUNT] = nuke
	actions[INSERT_CREDENTIAL] = insert
	actions[GET_CREDENTIAL] = get
	actions[RM_CREDENTIAL] = rm
	actions[GET_SERVICES] = services
	actions[UPDATE_CREDENTIAL] = update

	if len(argsWithoutProg) == 0 {
		log.Println("No action specified")
		help()
		return
	}

	if _, err := os.Stat(passmanHome); os.IsNotExist(err) {
		os.Create(passmanHome)
	}

	action, ok := actions[argsWithoutProg[0]]
	if ok {
		action()
	} else {
		log.Println("Invalid action specified")
		help()
	}
}

// func getUsernameAndPassword() (string, string) {
// 	fmt.Print("Username: ")
// 	text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
// 	text = cleanInput(text)
// 	fmt.Print("Password: ")
// 	bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
// 	password := cleanInput(string(bytePassword))
// 	fmt.Println("")
// 	return text, password
// }

func signin() {
	if len(argsWithoutProg) != 2 {
		log.Println("No account specified")
		return
	}

	username := argsWithoutProg[1]
	password := os.Getenv(PASSMAN_MASTER)

	var payload = `{"email":"%s","password":"%s","returnSecureToken": true}`
	payload = fmt.Sprintf(payload, username, password)
	req, err := http.NewRequest("GET", urlAuthUser, strings.NewReader(payload))

	if err != nil {
		log.Println(err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Println(err.Error())
	}

	if res.StatusCode != 200 {
		log.Println(string(body))
		return
	}

	log.Println("Login OK")

	err = ioutil.WriteFile(passmanHome, body, 0644)

	if err != nil {
		log.Println(err.Error())
	}
}

func cleanInput(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}

func getUserStore() ([]byte, error) {
	bytes, err := ioutil.ReadFile(passmanHome)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func genPassword() {
	if len(argsWithoutProg) != 2 {
		log.Println("length not specified")
		return
	}
	i, err := strconv.ParseInt(argsWithoutProg[1], 10, 32)
	if err != nil {
		fmt.Println(err.Error())
	}
	n, _ := GenerateRandomString(int(i))

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

func nuke() {

	tokenData, err := getUserStore()

	if err != nil {
		log.Println(err.Error())
		return
	}

	var storedJWT FirebaseStruct
	err = json.Unmarshal(tokenData, &storedJWT)

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("!!! WARNING !!!")
	log.Println("Are you sure you want to delete ALL data associaed with:", storedJWT.Email)
	log.Println("Type: 'yes' to remove all data")

	text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	text = cleanInput(text)
	if text != "yes" {
		log.Println("Aborting")
		return
	}

	var payload = `{"idToken":"%s"}`
	payload = fmt.Sprintf(payload, storedJWT.IDToken)

	// req, err := http.NewRequest("GET", urlAuthUser, strings.NewReader(payload))
	// fmt.Println(payload)
	req, err := http.NewRequest("DELETE", urlDeleteUser, strings.NewReader(payload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storedJWT.IDToken))

	if err != nil {
		log.Println(err.Error())
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	log.Println(string(body))
	os.Remove(passmanHome)
}

func register() {
	payloadCreateAccount := `{"email": "%s","password": "%s","returnSecureToken": true}`

	username := argsWithoutProg[1]
	password := os.Getenv(PASSMAN_MASTER)

	payloadCreateAccount = fmt.Sprintf(payloadCreateAccount, username, password)
	req, err := http.NewRequest("POST", urlCreateAccount, strings.NewReader(payloadCreateAccount))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println(string(body))
		return
	}

	err = ioutil.WriteFile(passmanHome, body, 0644)

	if err != nil {
		log.Println(err.Error())
		return
	}
	type response struct {
		IdToken string `json:"idToken"`
	}
	var id response
	json.Unmarshal(body, &id)

	payloadVerifyAccount := fmt.Sprintf(`{"requestType": "VERIFY_EMAIL","idToken": "%s"}`, id.IdToken)
	req, err = http.NewRequest("POST", urlVerifyAccount, strings.NewReader(payloadVerifyAccount))
	res, err = http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	body, _ = ioutil.ReadAll(res.Body)

	log.Println("Account created OK")
}

func insert() {
	if len(argsWithoutProg) != 4 {
		log.Println("Not enough args")
		return
	}

	newCredentialPayload := `{"serviceName": "%s","username": "%s","password": "%s"}`
	serviceName := ""
	tokenData, err := getUserStore()

	if err != nil {
		log.Println(err.Error())
		return
	}

	var storedJWT FirebaseStruct
	err = json.Unmarshal(tokenData, &storedJWT)

	if err != nil {
		log.Println(err.Error())
		return
	}

	username, password := argsWithoutProg[2], argsWithoutProg[3]

	password = encrypt([]byte(password), os.Getenv(PASSMAN_MASTER))
	username = encrypt([]byte(username), os.Getenv(PASSMAN_MASTER))
	serviceName = argsWithoutProg[1]
	newCredentialPayload = fmt.Sprintf(newCredentialPayload, serviceName, username, password)
	req, err := http.NewRequest("POST", urlNewCredential, strings.NewReader(newCredentialPayload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storedJWT.IDToken))

	if err != nil {
		log.Println(err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	body, _ := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println(string(body))
		log.Println("Try logging in again")
		return
	}

	log.Println("Credential added OK")
}

func update() {
	if len(argsWithoutProg) != 3 {
		log.Println("Not enough args")
		return
	}

	newCredentialPayload := `{"serviceName": "%s","password": "%s"}`

	tokenData, err := getUserStore()

	if err != nil {
		log.Println(err.Error())
		return
	}

	var storedJWT FirebaseStruct
	err = json.Unmarshal(tokenData, &storedJWT)

	if err != nil {
		log.Println(err.Error())
		return
	}

	username, password := argsWithoutProg[1], argsWithoutProg[2]
	password = encrypt([]byte(password), os.Getenv(PASSMAN_MASTER))
	newCredentialPayload = fmt.Sprintf(newCredentialPayload, username, password)
	req, err := http.NewRequest("POST", urlUpdateCredential, strings.NewReader(newCredentialPayload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storedJWT.IDToken))

	if err != nil {
		log.Println(err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	body, _ := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println(string(body))
		return
	}

	log.Println("Credential updated OK")
}

func get() {
	if len(argsWithoutProg) < 2 {
		log.Println("No service name")
		return
	}

	serviceName := ""
	tokenData, err := getUserStore()

	if err != nil {
		log.Println(err.Error())
		return
	}

	var storedJWT FirebaseStruct
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

	res, err := http.DefaultClient.Do(req)

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

	credentialRecord.Password = string(decrypt([]byte(sDec), os.Getenv(PASSMAN_MASTER)))
	credentialRecord.Username = string(decrypt([]byte(uName), os.Getenv(PASSMAN_MASTER)))

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

func services() {

	tokenData, err := getUserStore()

	if err != nil {
		log.Println(err.Error())
		return
	}

	var storedJWT FirebaseStruct
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

	res, err := http.DefaultClient.Do(req)

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
		ServiceName string
	}{}

	err = json.Unmarshal(body, &credentialRecord)

	if err != nil {
		log.Println(err.Error())
		return
	}

	data := [][]string{}
	for i := range credentialRecord {
		data = append(data, []string{credentialRecord[i].ServiceName})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Service"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func rm() {
	if len(argsWithoutProg) != 2 {
		log.Println("Not enough arguments")
		return
	}

	tokenData, err := getUserStore()

	if err != nil {
		log.Println(err.Error())
		return
	}

	var storedJWT FirebaseStruct
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

	res, err := http.DefaultClient.Do(req)

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

func encrypt(data []byte, passphrase string) string {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)

	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	sEnc := b64.StdEncoding.EncodeToString([]byte(ciphertext))
	return sEnc
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		panic(err.Error())
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		panic(err.Error())
	}

	return plaintext
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
