package db

import "gorm.io/gorm"

// QueryWithScan used to execute sql query like select
func QueryWithScan(db *gorm.DB, sql string, params []interface{}) ([]map[string]interface{}, error) {
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

// QueryWithExec used to execute sql query like insert, update, delete
func QueryWithExec(db *gorm.DB, sql string, params []interface{}) error {
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
