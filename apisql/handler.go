package apisql

import (
	"api-gateway-sql/logging"
	"api-gateway-sql/utils/httputil"

	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// ApiGetSqlHandler godoc
// @Summary      Get SQL Query result
// @Description  Trigger SQL query without params
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        target  path  string  true  "Target Name"
// @Success      200  {object}  httputil.HTTPResp
// @Failure      400  {object}  httputil.HTTPResp
// @Failure      500  {object}  httputil.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/{target} [get]
func (apisql *ApiSql) ApiGetSqlHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	targetName := vars["target"]

	target, database, err := getTargetAndDatabase(apisql, targetName)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	postParams := make(map[string]interface{}, 0)
	response, err := executeSingleSQLQuery(*target, *database, int(apisql.config.ApiGatewaySQL.Timeout), postParams)
	if err != nil {
		httputil.SendJSONResponse(resp, http.StatusInternalServerError, httputil.HTTPStatusInternalServerErrorMessage, nil)
		return
	}

	httputil.SendJSONResponse(resp, http.StatusOK, httputil.HTTPStatusOKMessage, response)
}

// ApiPostSqlHandler godoc
// @Summary      Get SQL Query result
// @Description  Trigger SQL query with params
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        target path  string  true  "Target Name"
// @Param        data  body  map[string]interface{}  true  "Data to send"
// @Success      200  {object}  httputil.HTTPResp
// @Failure      400  {object}  httputil.HTTPResp
// @Failure      500  {object}  httputil.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/{target} [post]
func (apisql *ApiSql) ApiPostSqlHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	targetName := vars["target"]

	target, database, err := getTargetAndDatabase(apisql, targetName)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	message := fmt.Sprintf("DBNAME : %s, TARGET NAME : %s", database.Dbname, target.Name)
	httputil.SendJSONResponse(resp, http.StatusOK, message, nil)
}

// InitializeDatabaseHandler godoc
// @Summary      Initialize Database
// @Description  Initialize Database by providing a sql query file
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        datasource  path  string  true  "Datasource Name"
// @Param        sqlfile  formData  file  true  "SQL Data to upload"
// @Success      200  {object}  httputil.HTTPResp
// @Failure      400  {object}  httputil.HTTPResp
// @Failure      500  {object}  httputil.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/{datasource}/init [post]
func (apisql *ApiSql) InitializeDatabaseHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	datasourceName := vars["datasource"]

	database, exist := apisql.config.GetDatabaseByDataSourceName(datasourceName)
	if !exist {
		err := fmt.Sprintf("the configured datasource name %s does not exist", datasourceName)
		httputil.SendJSONResponse(resp, http.StatusBadRequest, err, nil)
		return
	}

	file, _, err := req.FormFile("sqlfile")
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusBadRequest, "Unable to read the SQL file", nil)
		return
	}

	sqlBytes, err := io.ReadAll(file)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusBadRequest, "Unable to read the SQL file content", nil)
		return
	}

	err = executeInitSQLQuery(string(sqlBytes), database, int(apisql.config.ApiGatewaySQL.Timeout))
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusInternalServerError, "Unable to execute the SQL query", nil)
		return
	}

	httputil.SendJSONResponse(resp, http.StatusOK, httputil.HTTPStatusOKMessage, nil)
}
