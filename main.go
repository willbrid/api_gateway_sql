package main

import (
	"api-gateway-sql/apisql"
	"api-gateway-sql/config"
	"api-gateway-sql/db/stat"
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
// @license.url https://github.com/willbrid/easy_api_prom_sms_alert/blob/main/LICENSE
// @BasePath /v1
// @securityDefinitions.basic BasicAuth
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
		logging.Log(logging.Error, "error loading configuration : %s", err.Error())
		return
	}
	logging.Log(logging.Info, "configuration file '%s' was loaded successfully", configFile)

	// Execute automigration in database
	cnx, err := stat.Connect(configLoaded.ApiGatewaySQL.Sqlitedb)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return
	}
	cnx.AutoMigrate(&stat.BatchStatistic{}, &stat.FailureRange{})
	dbCnx, _ := cnx.DB()
	dbCnx.Close()

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

	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/api-gateway-sql/{target}", func(w http.ResponseWriter, r *http.Request) {
		apisql.ApiGetSqlHandler(w, r, *configLoaded)
	}).Methods("GET")
	v1.HandleFunc("/api-gateway-sql/{target}", func(w http.ResponseWriter, r *http.Request) {
		apisql.ApiPostSqlHandler(w, r, *configLoaded)
	}).Methods("POST")
	v1.HandleFunc("/api-gateway-sql/{target}/batch", func(w http.ResponseWriter, r *http.Request) {
		apisql.ApiPostSqlBatchHandler(w, r, *configLoaded)
	}).Methods("POST")
	v1.HandleFunc("/api-gateway-sql/{target}/batch", func(w http.ResponseWriter, r *http.Request) {
		apisql.ApiGetStatsHandler(w, r)
	}).Methods("GET")
	v1.HandleFunc("/api-gateway-sql/{datasource}/init", func(w http.ResponseWriter, r *http.Request) {
		apisql.InitializeDatabaseHandler(w, r, *configLoaded)
	}).Methods("POST")
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerUrl),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods("GET")
	router.Use(func(next http.Handler) http.Handler {
		return apisql.AuthMiddleware(next, *configLoaded)
	})

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
