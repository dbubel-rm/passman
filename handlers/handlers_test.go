package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/dbubel/passman/models"

	"github.com/dbubel/passman/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func fakeAuth(c *gin.Context) {}

func resetDB() *sqlx.DB {
	// db := fakeDB()
	cmdStr := "schema.sh"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	_, err := cmd.Output()

	db, err := sqlx.Connect("mysql", "root:@/passman_test")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

var a models.FirebaseAuthResp

func TestAddUser(t *testing.T) {
	db := resetDB()
	gin.SetMode(gin.ReleaseMode)
	testEngine := GetEngine(middleware.AuthUser, db)

	m := `{
			"email": "test@gmail.com2",
			"password": "123456k",
			returnSecureToken: true
	}`

	req, err := http.NewRequest("POST", "/user", strings.NewReader(m))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	testEngine.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	err = json.NewDecoder(resp.Body).Decode(&a)
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	db := resetDB()
	testEngine := GetEngine(middleware.AuthUser, db)
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(a)
	assert.NoError(t, err)
	// log.Println("sending to passman", b.String())
	req, err := http.NewRequest("DELETE", "/cred/user", strings.NewReader(b.String()))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	testEngine.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	var f models.FirebaseAuthResp
	json.NewDecoder(resp.Body).Decode(&f)

}
