package main

import (
	"db-api-sql/logging"

	"flag"

	"github.com/go-playground/validator/v10"
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

	// Load config file

	logging.Log(logging.Info, "configuration file '%s' was loaded successfully", configFile)
}
