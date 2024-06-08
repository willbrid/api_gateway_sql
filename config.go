package main

import (
	"time"
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
	Port     int    `mapstructure:"port" validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Dbname   string `mapstructure:"dbname" validate:"required"`
	Sslmode  bool   `mapstructure:"sslmode"`
}

type Target struct {
	Name           string   `mapstructure:"name" validate:"required,max=25"`
	DataSourceName string   `mapstructure:"data_source_name" validate:"required"`
	Datafields     []string `mapstructure:"datafields"`
	SqlQuery       string   `mapstructure:"sql" validate:"required"`
}

type Config struct {
	DbApiSQL struct {
		Timeout   time.Duration `mapstructure:"timeout" validate:"required"`
		Auth      `mapstructure:"auth"`
		Databases []Database `mapstructure:"databases" validate:"gt=0,required,dive"`
		Targets   []Target   `mapstructure:"targets" validate:"gt=0,required,dive"`
	} `mapstructure:"db_api_sql"`
}
