package main

import (
	"api-gateway-sql/apisql"
	"api-gateway-sql/config"
	_ "api-gateway-sql/docs"
	"api-gateway-sql/logging"

	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var validate *validator.Validate

// @title API GATEWAY SQL
// @version 1.0.0
// @description API used for executing SQL QUERY

// @contact.name API Support
// @contact.email ngaswilly77@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost
// @BasePath /
func main() {
	var (
		configFile  string
		listenPort  int
		enableHttps string
		certFile    string
		keyFile     string
	)

	flag.StringVar(&configFile, "config-file", "fixtures/config.default.yaml", "configuration file path")
	flag.StringVar(&certFile, "cert-file", "fixtures/tls/server.crt", "certificat file path")
	flag.StringVar(&keyFile, "key-file", "fixtures/tls/server.key", "private key file path")
	flag.StringVar(&enableHttps, "enable-https", "false", "configuration to enable https")
	flag.IntVar(&listenPort, "port", 5297, "listening port")
	flag.Parse()

	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Var(listenPort, "required,min=1024,max=49151"); err != nil {
		logging.Log(logging.Error, "you should provide a port number between 1024 and 49151")
		return
	}

	configLoaded, err := config.LoadConfig(configFile, validate)
	if err != nil {
		logging.Log(logging.Error, "error loading configuration")
		return
	}
	logging.Log(logging.Info, "configuration file '%s' was loaded successfully", configFile)

	strListenPort := strconv.Itoa(listenPort)
	boolEnableHttps, errParse := strconv.ParseBool(enableHttps)
	if errParse != nil {
		logging.Log(logging.Error, "unable to parse enable-https flag")
		return
	}

	swaggerUrl := fmt.Sprintf("http://localhost:%s/swagger/doc.json", strListenPort)
	if boolEnableHttps {
		swaggerUrl = fmt.Sprintf("https://localhost:%s/swagger/doc.json", strListenPort)
	}

	apisqlInstance := apisql.NewApiSql(configLoaded)
	router := mux.NewRouter()
	router.HandleFunc("/api-gateway-sql/{targetname}", apisqlInstance.ApiSqlHandler).Methods("GET", "POST")
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerUrl),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods("GET")
	router.Use(apisqlInstance.AuthMiddleware)

	logging.Log(logging.Info, "server is listening on port %v", strListenPort)
	if boolEnableHttps {
		logging.Log(logging.Info, "server is using https")
		err = http.ListenAndServeTLS(":"+strListenPort, certFile, keyFile, router)
	} else {
		logging.Log(logging.Info, "server is using http")
		err = http.ListenAndServe(":"+strListenPort, router)
	}

	if err != nil {
		logging.Log(logging.Error, "failed to start server: %v", err.Error())
		return
	}
}
