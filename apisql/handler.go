package apisql

import (
	"api-gateway-sql/utils/httputil"

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

	message := fmt.Sprintf("DBNAME : %s, TARGET NAME : %s", database.Dbname, target.Name)
	httputil.SendJSONResponse(resp, http.StatusOK, message, nil)
}
