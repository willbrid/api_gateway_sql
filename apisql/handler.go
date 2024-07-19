package apisql

import (
	"api-gateway-sql/logging"
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

	target, database, err := getTargetAndDatabase(apisql, targetName)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
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

	target, database, err := getTargetAndDatabase(apisql, targetName)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("DBNAME : %s, TARGET NAME : %s", database.Dbname, target.Name)
	httputil.SendJSONResponse(resp, http.StatusOK, message, nil)
}
