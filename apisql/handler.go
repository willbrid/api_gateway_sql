package apisql

import "net/http"

func (apisql *ApiSql) ApiSqlHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
}
