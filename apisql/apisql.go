package apisql

import "db-api-sql/config"

type ApiSql struct {
	config *config.Config
}

func NewApiSql(config *config.Config) *ApiSql {
	return &ApiSql{config}
}
