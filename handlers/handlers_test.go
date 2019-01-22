package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/dbubel/passman/middleware"
	"github.com/dbubel/passman/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func fakeAuth(c *gin.Context) {}

func resetDB(reset bool) *sqlx.DB {
	if reset {
		cmdStr := "sql_util.sh"
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		cmd.Output()
	}

	db, err := sqlx.Connect("mysql", "root:@/passman")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

var a models.FirebaseAuthResp

func TestAddUser(t *testing.T) {
	db := resetDB(true)
	// gin.SetMode(gin.ReleaseMode)
	testEngine := GetEngine(middleware.AuthUser, db)

	m := `{
			"email": "test@gmail.com2",
			"password": "123456k",
			returnSecureToken: true
	}`

	req, err := http.NewRequest("POST", "/public/user", strings.NewReader(m))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	testEngine.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	r, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(r, &a)
	assert.NoError(t, err)
	assert.Equal(t, "test@gmail.com2", a.Email)

}
func TestAddCredentials(t *testing.T) {
	db := resetDB(false)
	testEngine := GetEngine(middleware.AuthUser, db)

	payload := strings.NewReader(`{
		"serviceName":"gmail",
		"username":"username",
		"password":"password"
	}`)
	m := `Bearer %s`
	req, err := http.NewRequest("POST", "/credential", payload)
	req.Header.Set("Authorization", fmt.Sprintf(m, a.IDToken))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	testEngine.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

}

func TestGetCredential(t *testing.T) {
	db := resetDB(false)
	testEngine := GetEngine(middleware.AuthUser, db)

	// payload := strings.NewReader(`{
	// 	"serviceName":"gmail",
	// 	"username":"username",
	// 	"password":"password"
	// }`)
	m := `Bearer %s`
	req, err := http.NewRequest("GET", "/credential/gmail", nil)
	req.Header.Set("Authorization", fmt.Sprintf(m, a.IDToken))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	testEngine.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
func TestDeleteUser(t *testing.T) {
	db := resetDB(false)
	testEngine := GetEngine(middleware.AuthUser, db)

	m := `Bearer %s`

	req, err := http.NewRequest("DELETE", "/user", nil)
	req.Header.Set("Authorization", fmt.Sprintf(m, a.IDToken))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	testEngine.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
