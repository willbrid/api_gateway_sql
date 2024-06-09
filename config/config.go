package config

import (
	"db-api-sql/logging"

	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Auth struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username" validate:"required_if=Enabled true,min=2,max=25"`
	Password string `mapstructure:"password" validate:"required_if=Enabled true,min=8"`
}

type Database struct {
	Name     string `mapstructure:"name" validate:"required,max=25"`
	Type     string `mapstructure:"type" validate:"required,oneof=mysql mariadb postgresql mongodb"`
	Host     string `mapstructure:"host" validate:"required,ipv4"`
	Port     int    `mapstructure:"port" validate:"required,min=1024,max=49151"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Dbname   string `mapstructure:"dbname" validate:"required"`
	Sslmode  bool   `mapstructure:"sslmode"`
}

type Target struct {
	Name           string `mapstructure:"name" validate:"required,max=25"`
	DataSourceName string `mapstructure:"data_source_name" validate:"required"`
	Datafields     string `mapstructure:"datafields"`
	SqlQuery       string `mapstructure:"sql" validate:"required"`
}

type Config struct {
	DbApiSQL struct {
		Timeout   time.Duration `mapstructure:"timeout" validate:"required"`
		Auth      `mapstructure:"auth"`
		Databases []Database `mapstructure:"databases" validate:"gt=0,required,dive"`
		Targets   []Target   `mapstructure:"targets" validate:"gt=0,required,dive"`
	} `mapstructure:"db_api_sql"`
}

func setConfigDefaults(v *viper.Viper) {
	v.SetDefault("db_api_sql.timeout", "10s")
	v.SetDefault("db_api_sql.auth.enabled", false)
	v.SetDefault("db_api_sql.auth.username", "")
	v.SetDefault("db_api_sql.auth.password", "")
	v.SetDefault("db_api_sql.databases", make([]Database, 0))
	v.SetDefault("db_api_sql.targets", make([]Target, 0))
}

func LoadConfig(filename string, validate *validator.Validate) (*Config, error) {
	viper.SetConfigFile("yaml")
	viper.SetConfigFile(filename)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logging.Log(logging.Error, err.Error())
			return nil, err
		} else {
			logging.Log(logging.Error, err.Error())
			return nil, err
		}
	}

	var viperInstance *viper.Viper = viper.GetViper()
	// Set defaut configuration
	setConfigDefaults(viperInstance)

	// Parse configuration file to Config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}

	// Validate config struct
	if err := validate.Struct(config); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}

		for _, err := range err.(validator.ValidationErrors) {
			return nil, fmt.Errorf("validation failed on field '%s' for condition '%s'", err.Field(), err.Tag())
		}
	}

	return &config, nil
}
