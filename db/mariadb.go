package db

import (
	"api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MariadbDatabase struct {
	db *gorm.DB
}

func (mariadbDB *MariadbDatabase) Connect(db config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%vs", db.Username, db.Password, db.Host, db.Port, db.Dbname, db.Timeout)

	cnx, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err == nil {
		mariadbDB.db = cnx
	}

	return cnx, err
}

func (mariadbDB MariadbDatabase) ExecuteQuery(sqlQuery string, params map[string]interface{}) (SelectResult, error) {
	result, err := executeQuery(mariadbDB.db, sqlQuery, params)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (mariadbDB MariadbDatabase) ExecuteBatch(sqlQuery string, params []map[string]interface{}) error {
	return executeBatch(mariadbDB.db, sqlQuery, params)
}
