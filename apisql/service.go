package apisql

import (
	"api-gateway-sql/config"

	"fmt"
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
