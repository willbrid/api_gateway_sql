package db

import (
	"api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresInstance struct{}

func (i PostgresInstance) Connect(db config.Database, timeout int) (*gorm.DB, error) {
	sslMode := "disable"
	if db.Sslmode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%v connect_timeout=%v", db.Host, db.Username, db.Password, db.Dbname, db.Host, sslMode, timeout)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
}
