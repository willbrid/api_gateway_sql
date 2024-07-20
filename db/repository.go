package db

import (
	"strings"

	"gorm.io/gorm"
)

const (
	Select string = "SELECT"
)

// ExecuteWithScan used to execute sql query like select
func ExecuteWithScan(db *gorm.DB, sql string, params []interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	if len(params) > 0 {
		if err := db.Raw(sql, params...).Scan(&result).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Raw(sql).Scan(&result).Error; err != nil {
			return nil, err
		}
	}

	return result, nil
}

// ExecuteWithExec used to execute sql query like insert, update, delete
func ExecuteWithExec(db *gorm.DB, sql string, params []interface{}) error {
	if len(params) > 0 {
		if err := db.Exec(sql, params...).Error; err != nil {
			return err
		}
	} else {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetQuerySQLType used to get the base type (select, insert, update, delete,...) of a sql query
func GetSQLQueryType(sql string) string {
	return strings.ToUpper(strings.SplitN(sql, " ", 2)[0])
}
