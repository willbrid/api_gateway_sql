package db

import (
	"api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteDatabase struct {
	db *gorm.DB
}

func (sqliteDB *SqliteDatabase) Connect(db config.Database, timeout int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s.db", db.Dbname)

	cnx, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err == nil {
		sqliteDB.db = cnx
	}

	return cnx, err
}

func (sqliteDB SqliteDatabase) ExecuteQuery(query string, params []interface{}) (SelectResult, error) {
	result, err := executeQuery(sqliteDB.db, query, params)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (sqliteDB SqliteDatabase) ExecuteBatch(query string, batchSize int, bufferSize int) error {
	return nil
}
