package apisql

import (
	"api-gateway-sql/config"
	"api-gateway-sql/db"
	"api-gateway-sql/logging"

	"fmt"

	"gorm.io/gorm"
)

func getTargetAndDatabase(apisql *ApiSql, targetName string) (*config.Target, *config.Database, error) {
	target, exist := apisql.config.GetTargetByName(targetName)
	if !exist {
		return nil, nil, fmt.Errorf("the specified target name %s does not exist", targetName)
	}

	database, exist := apisql.config.GetDatabaseByDataSourceName(target.DataSourceName)
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

	dbInstance := db.GetDatabaseInstance(database.Type)
	cnx, err = dbInstance.Connect(database, timeout)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	sqlQueryType := db.GetSQLQueryType(target.SqlQuery)
	parsedQuery, params := db.TransformQuery(target.SqlQuery, postParams)

	switch sqlQueryType {
	case db.Select:
		result, err = db.ExecuteWithScan(cnx, parsedQuery, params)
		if err != nil {
			logging.Log(logging.Error, err.Error())
		}
	default:
		err = db.ExecuteWithExec(cnx, parsedQuery, params)
		if err != nil {
			logging.Log(logging.Error, err.Error())
		}
	}

	return result, err
}

func executeInitSQLQuery(sql string, database config.Database, timeout int) error {
	var (
		cnx *gorm.DB
		err error = nil
	)

	dbInstance := db.GetDatabaseInstance(database.Type)
	cnx, err = dbInstance.Connect(database, timeout)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	err = db.ExecuteWithExec(cnx, sql, nil)

	return err
}
