package utils

import (
	"crypto/rand"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
)

// I generated my own 8192 bit RSA keys, they are not verified by a 3rd party but are 100% secure.
func SkipTLS(r *http.Request) (*http.Response, error) {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	return client.Do(r)
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

func CleanInput(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}

func GetUserStore(passmanHome string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(passmanHome)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}
