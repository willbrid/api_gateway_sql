package db

import (
	"api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	db *gorm.DB
}

func (postgresDB *PostgresDatabase) Connect(db config.Database) (*gorm.DB, error) {
	sslMode := "disable"
	if db.Sslmode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%v connect_timeout=%v", db.Host, db.Username, db.Password, db.Dbname, db.Host, sslMode, int(db.Timeout))

	cnx, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err == nil {
		postgresDB.db = cnx
	}

	return cnx, err
}

func (postgresDB PostgresDatabase) ExecuteQuery(sqlQuery string, params map[string]interface{}) (SelectResult, error) {
	result, err := executeQuery(postgresDB.db, sqlQuery, params)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (postgresDB PostgresDatabase) ExecuteBatch(sqlQuery string, params []map[string]interface{}) error {
	return executeBatch(postgresDB.db, sqlQuery, params)
}
