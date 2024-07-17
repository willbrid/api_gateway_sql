package apisql

import (
	"api-gateway-sql/utils/httputil"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// ApiGetSqlHandler godoc
// @Summary      Get SQL Query result
// @Description  Trigger SQL query without params
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        targetname  path  string  true  "Target Name"
// @Success      200  {object}  httputil.HTTPResp
// @Failure      400  {object}  httputil.HTTPResp
// @Failure      404  {object}  httputil.HTTPResp
// @Failure      500  {object}  httputil.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/{targetname} [get]
func (apisql *ApiSql) ApiGetSqlHandler(resp http.ResponseWriter, req *http.Request) {
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

// ApiPostSqlHandler godoc
// @Summary      Get SQL Query result
// @Description  Trigger SQL query with params
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        targetname  path  string  true  "Target Name"
// @Param        data  body  map[string]interface{}  true  "Data to send"
// @Success      200  {object}  httputil.HTTPResp
// @Failure      400  {object}  httputil.HTTPResp
// @Failure      404  {object}  httputil.HTTPResp
// @Failure      500  {object}  httputil.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/{targetname} [post]
func (apisql *ApiSql) ApiPostSqlHandler(resp http.ResponseWriter, req *http.Request) {
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
