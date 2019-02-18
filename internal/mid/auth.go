package mid

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
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

		var respBody []byte
		resp, err := http.Get(PUBLIC_KEY_URL)

		if err != nil {
			return err
		}

		defer resp.Body.Close()

		respBody, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		err = json.Unmarshal(respBody, &publicPEM)
		if err != nil {
			return err
		}

		var publicKey *rsa.PublicKey
		var tok *jwt.Token
		// var err error

		token := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(token) < 2 {
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
			return errors.New("no valid token found")
		}

		iss, ok := tok.Claims.(jwt.MapClaims)["iss"].(string)
		if iss != JWT_ISSUER || !ok {
			return errors.New("Invalid ISS")
		}

		aud, ok := tok.Claims.(jwt.MapClaims)["aud"].(string)
		if aud != JWT_AUD || !ok {
			return errors.New("Invalid AUD")
		}

		localId, ok := tok.Claims.(jwt.MapClaims)["user_id"].(string)
		if !ok {
			return errors.New("no localId")
		}

		ctx := context.WithValue(r.Context(), "localId", localId)
		err = before(log, w, r.WithContext(ctx), params)
		return err
	}
}
