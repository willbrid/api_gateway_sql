package db

import (
	"api-gateway-sql/config"

	"gorm.io/gorm"
)

type DBInstance interface {
	Connect(dbConfig config.Database, timeout int) (*gorm.DB, error)
}

const (
	Mariadb    string = "mariadb"
	MySQL      string = "mysql"
	PostgreSQL string = "postgres"
	Sqlserver  string = "sqlserver"
	SQLite     string = "sqlite"
)

func GetDatabaseInstance(dbType string) DBInstance {
	switch dbType {
	case Mariadb:
	case MySQL:
		return MySQLInstance{}
	case PostgreSQL:
		return PostgresInstance{}
	case Sqlserver:
		return SqlserverInstance{}
	case SQLite:
		return SqliteInstance{}
	}

	return SqliteInstance{}
}
