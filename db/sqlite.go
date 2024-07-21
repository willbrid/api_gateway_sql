package db

import (
	"api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteInstance struct{}

func (i SqliteInstance) Connect(db config.Database, timeout int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s.db", db.Dbname)

	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
}
