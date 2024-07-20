package db

import (
	"api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLInstance struct{}

func (i MySQLInstance) Connect(db config.Database, timeout int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%vs", db.Username, db.Password, db.Host, db.Port, db.Dbname, timeout)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
