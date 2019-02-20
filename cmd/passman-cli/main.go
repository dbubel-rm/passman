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

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
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
var passmanHome = "~/.passman/session.json"

var argsWithoutProg = os.Args[1:]

const (
	CREATE_ACCOUNT = "create:account"
	DELETE_ACCOUNT = "delete:account"
	HELP           = "help"
	VERSION        = "version"
	GEN_PASS       = "gen:pass"
	ADD_CRED       = "add:credential"
	LOGIN          = "login"
	PASSMAN_MASTER = "PASSMAN_MASTER"
	GET_CRED       = "get:credential"
	GET_CREDS      = "get:credentials"
	DELETE_CRED    = "delete:credential"
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

func displayOptions() {
	fmt.Println("Passman is a utility for managing your passwords.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\tpassman <command> [arguments]")
	fmt.Println("")
	fmt.Println("The commands are:")
	fmt.Println("")
	fmt.Printf("\t%s\t\tCreates a new passman account. Ex) passman %s newexample@example.com\n", CREATE_ACCOUNT, CREATE_ACCOUNT)
	fmt.Printf("\t%s\t\tDeletes ALL credentials saved under you active account. Ex) passman %s\n", DELETE_ACCOUNT, DELETE_ACCOUNT)
	fmt.Printf("\t%s\t\tAdd a credential. Ex) passman %s serviceName\n", ADD_CRED, ADD_CRED)
	fmt.Printf("\t%s\t\tGet a stored credential. Ex) passman %s serviceName\n", GET_CRED, GET_CRED)
	fmt.Printf("\t%s\t\tGet all stored credentials. Ex) passman %s\n", GET_CREDS, GET_CREDS)
	fmt.Printf("\t%s\t\tDeletes a stored credential. Ex) passman %s service_name\n", DELETE_CRED, DELETE_CRED)
	fmt.Printf("\t%s\t\tAuthenticate a passman session good for 30 minutes\n", LOGIN)
	fmt.Printf("\t%s\t\tGenerates a crypto random string to be used for a secure password\n", GEN_PASS)
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
		log.Println("No PASSMAN_MASTER environment vairable set")
		return
	}

	log.SetFlags(log.Lshortfile)

	actions := make(map[string]func())
	actions[HELP] = displayOptions
	actions[VERSION] = version
	actions[GEN_PASS] = genPassword
	// API calls
	actions[LOGIN] = signin
	actions[CREATE_ACCOUNT] = createAccount
	actions[DELETE_ACCOUNT] = deleteAccount
	actions[ADD_CRED] = addCredential
	actions[GET_CRED] = getCredential
	actions[DELETE_CRED] = deleteCredential
	// actions[GET_CREDS] = getCredentials

	if len(argsWithoutProg) == 0 {
		log.Println("No action specified")
		displayOptions()
		return
	}

	if _, err := os.Stat(passmanHome); os.IsNotExist(err) {
		os.Create(passmanHome)
	}

	f, ok := actions[argsWithoutProg[0]]
	if ok {
		f()
	} else {
		log.Println("Invalid action specified")
		displayOptions()
	}
}

func getUsernameAndPassword() (string, string) {
	fmt.Print("Username: ")
	text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	text = cleanInput(text)
	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	password := cleanInput(string(bytePassword))
	fmt.Println("")
	return text, password
}

func signin() {
	if len(argsWithoutProg) < 2 {
		log.Println("No account email specified")
		return
	}
	username := argsWithoutProg[1]
	password := os.Getenv(PASSMAN_MASTER)

	var payload = `{"email":"%s","password":"%s","returnSecureToken": true}`
	payload = fmt.Sprintf(payload, username, password)
	fmt.Println(payload)
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

	err = ioutil.WriteFile(passmanHome, body, 0644)

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Login OK")
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
	fmt.Println(n)
}

func deleteAccount() {

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

	if err != nil {
		log.Println(err.Error())
	}

	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storedJWT.IDToken))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	log.Println(string(body))
	os.Remove(passmanHome)
}

func createAccount() {
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

	log.Println("Account created OK")
}

func addCredential() {
	if len(argsWithoutProg) < 2 {
		log.Println("No service name")
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

	username, password := getUsernameAndPassword()
	fmt.Print("Confirm: ")

	bytePasswordConfirm, err := terminal.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	passwordConfirm := cleanInput(string(bytePasswordConfirm))
	fmt.Println()

	if password != passwordConfirm {
		fmt.Println("Passwords do not match")
		return
	}

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
		return
	}

	log.Println("Credential added OK")
}

func getCredential() {
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
	log.Println("Credential retreived OK:", credentialRecord)
}

// func getCredentials() {
// 	tokenData, err := getUserStore()

// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}

// 	var storedJWT FirebaseStruct
// 	err = json.Unmarshal(tokenData, &storedJWT)

// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	req, err := http.NewRequest("GET", urlNewCredential+"s", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storedJWT.IDToken))

// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}

// 	res, err := http.DefaultClient.Do(req)

// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}

// 	body, _ := ioutil.ReadAll(res.Body)

// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	defer res.Body.Close()

// 	if res.StatusCode != 200 {
// 		log.Println(string(body))
// 		return
// 	}

// 	var credentialRecord = []struct {
// 		ServiceName string
// 		Username    string
// 		Password    string
// 	}{}

// 	err = json.Unmarshal(body, &credentialRecord)

// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}

// 	for i := 0; i < len(credentialRecord); i++ {
// 		sDec, err := b64.StdEncoding.DecodeString(credentialRecord[i].Password)

// 		if err != nil {
// 			log.Println(err.Error())
// 			return
// 		}

// 		credentialRecord[i].Password = string(decrypt([]byte(sDec), os.Getenv(PASSMAN_MASTER)))
// 	}

// 	log.Println("Credential retreived OK:", credentialRecord)
// }

func deleteCredential() {
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

	// var credentialRecord = struct {
	// 	ServiceName string
	// 	Username    string
	// 	Password    string
	// }{}

	// err = json.Unmarshal(body, &credentialRecord)

	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }

	// sDec, err := b64.StdEncoding.DecodeString(credentialRecord.Password)

	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }

	// credentialRecord.Password = string(decrypt([]byte(sDec), os.Getenv(PASSMAN_MASTER)))
	log.Println("Credential deleted OK:", string(body))
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
