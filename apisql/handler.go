package apisql

import (
	"api-gateway-sql/config"
	"api-gateway-sql/db"
	"api-gateway-sql/db/stat"
	"api-gateway-sql/logging"
	"api-gateway-sql/utils/httputil"

	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

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
func ApiGetSqlHandler(resp http.ResponseWriter, req *http.Request, configLoaded config.Config) {
	var (
		vars       map[string]string = mux.Vars(req)
		targetName string            = vars["target"]
		err        error
		target     *config.Target
		database   *config.Database
		response   db.SelectResult
	)

	target, database, err = getTargetAndDatabase(configLoaded, targetName)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	postParams := make(map[string]interface{}, 0)
	response, err = executeSingleSQLQuery(*target, *database, postParams)
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
func ApiPostSqlHandler(resp http.ResponseWriter, req *http.Request, configLoaded config.Config) {
	var (
		vars       map[string]string = mux.Vars(req)
		targetName string            = vars["target"]
		err        error
		target     *config.Target
		database   *config.Database
		response   db.SelectResult
	)

	target, database, err = getTargetAndDatabase(configLoaded, targetName)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var postParams map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&postParams); err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusBadGateway, err.Error(), nil)
		return
	}

	response, err = executeSingleSQLQuery(*target, *database, postParams)
	if err != nil {
		httputil.SendJSONResponse(resp, http.StatusInternalServerError, httputil.HTTPStatusInternalServerErrorMessage, nil)
		return
	}

	httputil.SendJSONResponse(resp, http.StatusOK, httputil.HTTPStatusOKMessage, response)
}

// ApiPostSqlBatchHandler godoc
// @Summary      Execute batch sql query
// @Description  Execute batch sql query with values from a csv file
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        target path  string  true  "Target Name"
// @Param        csvfile  formData  file  true  "CSV Data to import"
// @Success      200  {object}  httputil.HTTPResp
// @Failure      400  {object}  httputil.HTTPResp
// @Failure      500  {object}  httputil.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/{target}/batch [post]
func ApiPostSqlBatchHandler(resp http.ResponseWriter, req *http.Request, configLoaded config.Config) {
	var (
		vars       map[string]string = mux.Vars(req)
		targetName string            = vars["target"]
		err        error
		postFile   multipart.File
		target     *config.Target
		database   *config.Database
	)

	target, database, err = getTargetAndDatabase(configLoaded, targetName)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	postFile, _, err = req.FormFile("csvfile")
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusBadRequest, "Unable to read the SQL file", nil)
		return
	}

	go func() {
		err = executeBatchSQLQuery(configLoaded.ApiGatewaySQL.Sqlitedb, *target, *database, postFile)
		if err != nil {
			logging.Log(logging.Error, err.Error())
			return
		}
	}()

	httputil.SendJSONResponse(resp, http.StatusOK, httputil.HTTPStatusOKMessage, nil)
}

// ApiGetStatsHandler godoc
// @Summary      Get statistics
// @Description  Get all batch statistics
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        page_num query  int  false  "Page number" default(1)
// @Param        page_size query int  false  "Page size" default(20)
// @Success      200  {object}  httputil.HTTPResp
// @Failure      400  {object}  httputil.HTTPResp
// @Failure      500  {object}  httputil.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/stats [get]
func ApiGetStatsHandler(resp http.ResponseWriter, req *http.Request, configLoaded config.Config) {
	var (
		queries  url.Values = req.URL.Query()
		pageNum  int
		pageSize int
		err      error
		stats    []stat.BatchStatistic
	)

	pageNum, err = strconv.Atoi(queries.Get("page_num"))
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusBadRequest, "Unable to handle page_num", nil)
		return
	}
	pageSize, err = strconv.Atoi(queries.Get("page_size"))
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusBadRequest, "Unable to handle page_size", nil)
		return
	}

	stats, err = getStats(configLoaded.ApiGatewaySQL.Sqlitedb, pageNum, pageSize)
	if err != nil {
		httputil.SendJSONResponse(resp, http.StatusInternalServerError, httputil.HTTPStatusInternalServerErrorMessage, nil)
		return
	}

	httputil.SendJSONResponse(resp, http.StatusOK, httputil.HTTPStatusOKMessage, stats)
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
func InitializeDatabaseHandler(resp http.ResponseWriter, req *http.Request, configLoaded config.Config) {
	var (
		vars           map[string]string = mux.Vars(req)
		datasourceName string            = vars["datasource"]
		err            error
		database       config.Database
		exist          bool
		file           multipart.File
		sqlBytes       []byte
	)

	database, exist = configLoaded.GetDatabaseByDataSourceName(datasourceName)
	if !exist {
		err := fmt.Sprintf("the configured datasource name %s does not exist", datasourceName)
		httputil.SendJSONResponse(resp, http.StatusBadRequest, err, nil)
		return
	}

	file, _, err = req.FormFile("sqlfile")
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusBadRequest, "Unable to read the SQL file", nil)
		return
	}

	sqlBytes, err = io.ReadAll(file)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusBadRequest, "Unable to read the SQL file content", nil)
		return
	}

	err = executeInitSQLQuery(string(sqlBytes), database)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		httputil.SendJSONResponse(resp, http.StatusInternalServerError, "Unable to execute the SQL query", nil)
		return
	}

	httputil.SendJSONResponse(resp, http.StatusOK, httputil.HTTPStatusOKMessage, nil)
}
