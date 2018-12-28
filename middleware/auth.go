package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	JWT_ISSUER     string = "https://securetoken.google.com/passman-fc9e0"
	JWT_AUD        string = "passman-fc9e0"
	PUBLIC_KEY_URL string = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
)

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
}

func AuthUser(c *gin.Context) {
	var publicKey *rsa.PublicKey
	var tok *jwt.Token
	var err error

	token := strings.Split(c.GetHeader("Authorization"), "Bearer ")

	if len(token) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auth header"})
		c.Abort()
		return
	}

	s := strings.Replace(token[1], " ", "", -1)

	for _, pem := range publicPEM {
		tok, err = jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       err.Error(),
			"description": "Error validating jwt",
		})
		return
	}

	iss, ok := tok.Claims.(jwt.MapClaims)["iss"].(string)
	if iss != JWT_ISSUER || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ISS"})
		c.Abort()
		return
	}

	aud, ok := tok.Claims.(jwt.MapClaims)["aud"].(string)
	if aud != JWT_AUD || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid AUD"})
		c.Abort()
		return
	}
	// c.Set("claimsMap", tok.Claims.(jwt.MapClaims))
	log.Println(tok.Claims.(jwt.MapClaims))
	c.Set("userID", tok.Claims.(jwt.MapClaims)["user_id"])
	c.Set("email", tok.Claims.(jwt.MapClaims)["email"])
	c.Set("jwt", s)
	c.Next()
}
