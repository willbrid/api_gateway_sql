package apisql

import "api-gateway-sql/config"

type ApiSql struct {
	config *config.Config
}

func NewApiSql(config *config.Config) *ApiSql {
	return &ApiSql{config}
}
