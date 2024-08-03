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

func (postgresDB *PostgresDatabase) Connect(db config.Database, timeout int) (*gorm.DB, error) {
	sslMode := "disable"
	if db.Sslmode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%v connect_timeout=%v", db.Host, db.Username, db.Password, db.Dbname, db.Host, sslMode, timeout)

	cnx, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err == nil {
		postgresDB.db = cnx
	}

	return cnx, err
}

func (postgresDB PostgresDatabase) ExecuteQuery(query string, params []interface{}) (SelectResult, error) {
	result, err := executeQuery(postgresDB.db, query, params)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (postgresDB PostgresDatabase) ExecuteBatch(query string, batchSize int, bufferSize int) error {
	return nil
}
