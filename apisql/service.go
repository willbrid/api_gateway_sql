package apisql

import (
	"api-gateway-sql/config"
	"api-gateway-sql/db"
	"api-gateway-sql/logging"

	"fmt"

	"gorm.io/gorm"
)

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

	parsedQuery, params := db.TransformQuery(target.SqlQuery, postParams)

	result, err = dbInstance.ExecuteQuery(parsedQuery, params)

	return result, err
}

func executeInitSQLQuery(sql string, database config.Database, timeout int) error {
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

	_, err = dbInstance.ExecuteQuery(sql, make([]interface{}, 0))

	return err
}
