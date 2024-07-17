package apisql

import (
	"api-gateway-sql/config"
	"api-gateway-sql/utils/file"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

const configContent string = `---
api_gateway_sql:
  timeout: "10s"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "sqlite"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
  targets:
  - name: xxxxx
    data_source_name: xxxxx
    datafields: ""
    sql: "select * from students"
`

func triggerTest(t *testing.T, statusCode int, credential string) {
	filename, err := file.CreateConfigFileForTesting(configContent)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.Remove(filename)

	configLoaded, err := config.LoadConfig(filename, validate)
	if err != nil {
		t.Fatalf(err.Error())
	}

	apisqlInstance := NewApiSql(configLoaded)

	req, err := http.NewRequest("GET", "/api-gateway-sql/xxxxx", nil)
	if err != nil {
		t.Fatal(err)
	}

	if credential != "" {
		req.Header.Add("Authorization", credential)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api-gateway-sql/{targetname}", apisqlInstance.ApiSqlHandler).Methods("GET")
	router.Use(apisqlInstance.AuthMiddleware)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != statusCode {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, statusCode)
	}
}

func TestNoAuthorizationHeaderAuthMiddleware(t *testing.T) {
	triggerTest(t, http.StatusUnauthorized, "")
}

func TestInvalidAuthorizationHeaderAuthMiddleware(t *testing.T) {
	triggerTest(t, http.StatusUnauthorized, "xxxxx")
}

func TestFailedToDecodeBase64TokenAuthMiddleware(t *testing.T) {
	triggerTest(t, http.StatusUnauthorized, "Basic xxxxx")
}

func TestInvalidUsernameOrPasswordAuthMiddleware(t *testing.T) {
	triggerTest(t, http.StatusUnauthorized, "Basic eHh4eHg6eHh4")
}

func TestCorrectUsernameAndPasswordAuthMiddleware(t *testing.T) {
	triggerTest(t, http.StatusOK, "Basic eHh4eHg6eHh4eHh4eHg=")
}
