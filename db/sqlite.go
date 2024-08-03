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

func (sqliteDB SqliteDatabase) ExecuteQuery(sqlQuery string, params map[string]interface{}) (SelectResult, error) {
	result, err := executeQuery(sqliteDB.db, sqlQuery, params)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (sqliteDB SqliteDatabase) ExecuteBatch(sqlQuery string, params []map[string]interface{}) error {
	return executeBatch(sqliteDB.db, sqlQuery, params)
}
