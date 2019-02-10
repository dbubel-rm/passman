package mid

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dbubel/passman/internal/platform/web"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

const (
	JWT_ISSUER     string = "https://securetoken.google.com/passman-fc9e0"
	JWT_AUD        string = "passman-fc9e0"
	PUBLIC_KEY_URL string = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
)

var publicPEM map[string]string

// RequestLogger writes some information about the request to the logs in
// the format: TraceID : (200) GET /foo -> IP ADDR (latency)
func AuthHandler(before web.Handler) web.Handler {
	// Wrap this handler around the next one provided.
	return func(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
		fmt.Println("AUTH")
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

		log.Println(respBody)

		var publicKey *rsa.PublicKey
		var tok *jwt.Token
		// var err error

		token := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(token) != 2 {
			// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auth header"})
			// c.Abort()
			return errors.New("no token")
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
			// c.Abort()
			// c.JSON(http.StatusInternalServerError, gin.H{
			// 	"error":       err.Error(),
			// 	"description": "Error validating jwt",
			// })
			// return
		}

		iss, ok := tok.Claims.(jwt.MapClaims)["iss"].(string)
		if iss != JWT_ISSUER || !ok {
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ISS"})
			// c.Abort()
			// return
		}
		log.Println("AUth OK")

		aud, ok := tok.Claims.(jwt.MapClaims)["aud"].(string)
		if aud != JWT_AUD || !ok {
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid AUD"})
			// c.Abort()
			// return
		}

		ctx := context.WithValue(r.Context(), "Username", "HI")
		err = before(log, w, r.WithContext(ctx), params)
		// log.Printf("%s -> %d -> %s -> %s", r.Method, r.ContentLength, r.URL.Path, r.RemoteAddr)
		// For consistency return the error we received.
		return err
	}
}
