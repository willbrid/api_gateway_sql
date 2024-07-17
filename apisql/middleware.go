package apisql

import (
	"api-gateway-sql/logging"

	"encoding/base64"
	"net/http"
	"strings"
)

func (apiSql *ApiSql) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var auth string = req.Header.Get("Authorization")

		if apiSql.config.ApiGatewaySQL.Auth.Enabled && !strings.HasPrefix(req.RequestURI, "/swagger/") {
			if auth == "" {
				logging.Log(logging.Error, "no authorization header found")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(auth, "Basic ") {
				logging.Log(logging.Error, "invalid authorization header")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Basic ")
			decodedToken, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				logging.Log(logging.Error, "failed to decode base64 token - %v", err)
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			credentialParts := strings.SplitN(string(decodedToken), ":", 2)
			username := credentialParts[0]
			password := credentialParts[1]
			if username != apiSql.config.ApiGatewaySQL.Auth.Username || password != apiSql.config.ApiGatewaySQL.Auth.Password {
				logging.Log(logging.Error, "invalid username or password")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(resp, req)
	})
}
