package apisql

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (apisql *ApiSql) ApiSqlHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	targetName := vars["targetname"]

	target, exist := apisql.config.GetTargetByName(targetName)
	if !exist {
		http.Error(resp, fmt.Sprintf("the specified target name %s does not exist", targetName), http.StatusBadRequest)
		return
	}

	database, exist := apisql.config.GetDatabaseByDataSourceName(target.DataSourceName)
	if !exist {
		http.Error(resp, fmt.Sprintf("the configured datasource name %s does not exist", target.DataSourceName), http.StatusInternalServerError)
		return
	}

	fmt.Println("DBNAME : "+database.Dbname, " TARGET NAME : "+target.Name)

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
}
