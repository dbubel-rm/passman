package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	JWT_ISSUER     string = "https://securetoken.google.com/passman-fc9e0"
	JWT_AUD        string = "passman-fc9e0"
	PUBLIC_KEY_URL string = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
)

type firebaseAuthResp struct {
	Kind         string `json:"kind"`
	LocalID      string `json:"localId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	IDToken      string `json:"idToken"`
	Registered   bool   `json:"registered"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

var publicPEM map[string]string

func init() {
	var respBody []byte
	resp, err := http.Get(PUBLIC_KEY_URL)

	if err != nil {
		log.Println(err.Error())
	}

	defer resp.Body.Close()

	respBody, _ = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
	}

	err = json.Unmarshal(respBody, &publicPEM)
	if err != nil {
		log.Println(err.Error())
	}
	// fmt.Println(publicPEM)
}

func AuthUser(c *gin.Context) {
	var publicKey *rsa.PublicKey
	var tok *jwt.Token
	var err error
	var a firebaseAuthResp

	if err = c.BindJSON(&a); err != nil {
		log.Println(err.Error())
		c.Abort()
		return
	}

	for _, pem := range publicPEM {
		tok, err = jwt.Parse(a.IDToken, func(token *jwt.Token) (interface{}, error) {
			publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			return publicKey, err
		})

		if err != nil {
			continue
		} else {
			break
		}
	}

	// No valid jwt was found
	if err != nil {
		c.Abort()
		return
	}

	iss, ok := tok.Claims.(jwt.MapClaims)["iss"].(string)
	if iss != JWT_ISSUER || !ok {
		log.Println("Invalid iss claim")
		c.Abort()
		return
	}

	aud, ok := tok.Claims.(jwt.MapClaims)["aud"].(string)
	if aud != JWT_AUD || !ok {
		log.Println("Invalid aud claim")
		c.Abort()
		return
	}
}
