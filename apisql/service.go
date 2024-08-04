package apisql

import (
	"api-gateway-sql/config"
	"api-gateway-sql/db"
	"api-gateway-sql/db/stat"
	"api-gateway-sql/logging"

	"fmt"

	"gorm.io/gorm"
)

// getTargetAndDatabase used to get target and his database
func getTargetAndDatabase(config config.Config, targetName string) (*config.Target, *config.Database, error) {
	target, exist := config.GetTargetByName(targetName)
	if !exist {
		return nil, nil, fmt.Errorf("the specified target name %s does not exist", targetName)
	}

	database, exist := config.GetDatabaseByDataSourceName(target.DataSourceName)
	if !exist {
		return nil, nil, fmt.Errorf("the configured datasource name %s does not exist", target.DataSourceName)
	}

	return &target, &database, nil
}

// executeSingleSQLQuery used to execute a single sql query from target
func executeSingleSQLQuery(target config.Target, database config.Database, timeout int, postParams map[string]interface{}) (db.SelectResult, error) {
	var (
		cnx    *gorm.DB
		result db.SelectResult = nil
		err    error           = nil
	)

	dbInstance := db.NewDatabase(database.Type)
	cnx, err = dbInstance.Connect(database, timeout)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	result, err = dbInstance.ExecuteQuery(target.SqlQuery, postParams)

	return result, err
}

// executeBatchSQLQuery used to execute a batch query from target
func executeBatchSQLQuery(target config.Target, database config.Database, timeout int, filePath string) error {
	var (
		cnx *gorm.DB
		// stat stat.BatchStatistic = stat.NewBatchStatistic(target.Name)
		err error = nil
	)

	dbInstance := db.NewDatabase(database.Type)
	cnx, err = dbInstance.Connect(database, timeout)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	return nil
}

// getStats used to get stats in application database
func getStats(sqlitedb string, pageNum int, pageSize int) ([]stat.BatchStatistic, error) {
	var (
		defaultPageNum  int = 1
		defaultPageSize int = 20
	)

	if pageNum > 0 && pageSize > 0 {
		defaultPageNum = pageNum
		defaultPageSize = pageSize
	}

	return stat.GetBatchStatistics(sqlitedb, defaultPageNum, defaultPageSize)
}

// mapBatchFieldToValueLine used to construct a record of a batch sql query
func mapBatchFieldToValueLine(fields []string, values []string) (map[string]interface{}, error) {
	if len(fields) != len(values) {
		return nil, fmt.Errorf("bad mapping fields and file column")
	}

	var result map[string]interface{}

	for index, field := range fields {
		result[field] = values[index]
	}

	return result, nil
}

// executeInitSQLQuery used to populate a target database for testing
func executeInitSQLQuery(sqlQuery string, database config.Database, timeout int) error {
	var (
		cnx *gorm.DB
		err error = nil
	)

	dbInstance := db.NewDatabase(database.Type)
	cnx, err = dbInstance.Connect(database, timeout)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	_, err = dbInstance.ExecuteQuery(sqlQuery, make(map[string]interface{}, 0))

	return err
}
