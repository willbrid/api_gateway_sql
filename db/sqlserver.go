package db

import (
	"api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type SqlserverInstance struct{}

func (i SqlserverInstance) Connect(db config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%v?database=%s", db.Username, db.Password, db.Host, db.Port, db.Dbname)

	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
}
