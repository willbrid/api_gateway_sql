package main

import (
	"db-api-sql/apisql"
	"db-api-sql/config"
	"db-api-sql/logging"

	"flag"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

func main() {
	var (
		configFile string
		listenPort int
	)

	flag.StringVar(&configFile, "config-file", "config.default.yaml", "configuration file path")
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

	apisqlInstance := apisql.NewApiSql(configLoaded)
	router := mux.NewRouter()
	router.HandleFunc("/dbapisql/{targetname}", apisqlInstance.ApiSqlHandler).Methods("POST")

	strListenPort := strconv.Itoa(listenPort)
	logging.Log(logging.Info, "server is listening on port %v", strListenPort)
	err = http.ListenAndServe(":"+strListenPort, router)
	if err != nil {
		logging.Log(logging.Error, "failed to start server: %v", err.Error())
		return
	}
}
