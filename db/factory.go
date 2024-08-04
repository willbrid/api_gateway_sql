package db

import (
	"api-gateway-sql/config"

	"gorm.io/gorm"
)

type IDatabase interface {
	Connect(dbConfig config.Database) (*gorm.DB, error)
	ExecuteQuery(sqlQuery string, params map[string]interface{}) (SelectResult, error)
	ExecuteBatch(sqlQuery string, params []map[string]interface{}) error
}

const (
	Mariadb    string = "mariadb"
	MySQL      string = "mysql"
	PostgreSQL string = "postgres"
	Sqlserver  string = "sqlserver"
	SQLite     string = "sqlite"
)

func NewDatabase(dbType string) IDatabase {
	switch dbType {
	case Mariadb:
		return &MariadbDatabase{}
	case MySQL:
		return &MySQLDatabase{}
	case PostgreSQL:
		return &PostgresDatabase{}
	case Sqlserver:
		return &SqlserverDatabase{}
	case SQLite:
		return &SqliteDatabase{}
	default:
		return &SqliteDatabase{}
	}
}
