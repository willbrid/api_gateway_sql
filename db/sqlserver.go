package db

import (
	"api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type SqlserverDatabase struct {
	db *gorm.DB
}

func (sqlserverDB *SqlserverDatabase) Connect(db config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%v?database=%s&connection+timeout=%v", db.Username, db.Password, db.Host, db.Port, db.Dbname, db.Timeout)

	cnx, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err == nil {
		sqlserverDB.db = cnx
	}

	return cnx, err
}

func (sqlserverDB SqlserverDatabase) ExecuteQuery(sqlQuery string, params map[string]interface{}) (SelectResult, error) {
	result, err := executeQuery(sqlserverDB.db, sqlQuery, params)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (sqlserverDB SqlserverDatabase) ExecuteBatch(sqlQuery string, params []map[string]interface{}) error {
	return executeBatch(sqlserverDB.db, sqlQuery, params)
}
