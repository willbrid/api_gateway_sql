package db

import (
	"api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDatabase struct {
	db *gorm.DB
}

func (mysqlDB *MySQLDatabase) Connect(db config.Database, timeout int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%vs", db.Username, db.Password, db.Host, db.Port, db.Dbname, timeout)

	cnx, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err == nil {
		mysqlDB.db = cnx
	}

	return cnx, err
}

func (mysqlDB MySQLDatabase) ExecuteQuery(query string, params []interface{}) (SelectResult, error) {
	result, err := executeQuery(mysqlDB.db, query, params)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (mysqlDB MySQLDatabase) ExecuteBatch(query string, batchSize int, bufferSize int) error {
	return nil
}
