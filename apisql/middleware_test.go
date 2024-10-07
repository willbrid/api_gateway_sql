package apisql_test

import (
	"api-gateway-sql/apisql"
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
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "sqlite"
    dbname: "/tmp/xxxxx"
    timeout: "10s"
  targets:
  - name: xxxxx
    data_source_name: xxxxx
    datafields: ""
    sql: "select * from students"
`

func triggerTest(t *testing.T, statusCode int, credential string) {
	filename, err := file.CreateConfigFileForTesting(configContent)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(filename)

	configLoaded, err := config.LoadConfig(filename, validate)
	if err != nil {
		t.Fatal(err.Error())
	}

	req, err := http.NewRequest("GET", "/api-gateway-sql/xxxxx", nil)
	if err != nil {
		t.Fatal(err)
	}

	if credential != "" {
		req.Header.Add("Authorization", credential)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api-gateway-sql/{targetname}", func(w http.ResponseWriter, r *http.Request) {
		apisql.ApiGetSqlHandler(w, r, *configLoaded)
	}).Methods("GET")
	router.Use(func(next http.Handler) http.Handler {
		return apisql.AuthMiddleware(next, *configLoaded)
	})
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != statusCode {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, statusCode)
	}
}

func TestAuthentication(t *testing.T) {
	t.Run("No authorization header", func(t *testing.T) {
		triggerTest(t, http.StatusUnauthorized, "")
	})

	t.Run("Invalid authorization header", func(t *testing.T) {
		triggerTest(t, http.StatusUnauthorized, "xxxxx")
	})

	t.Run("Failed to decode base64 token", func(t *testing.T) {
		triggerTest(t, http.StatusUnauthorized, "Basic xxxxx")
	})

	t.Run("Invalid username or password", func(t *testing.T) {
		triggerTest(t, http.StatusUnauthorized, "Basic eHh4eHg6eHh4")
	})
}
