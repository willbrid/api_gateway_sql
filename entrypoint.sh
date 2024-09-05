#!/bin/sh

exec /usr/local/bin/api-gateway-sql \
  --config-file /etc/api-gateway-sql/config.yaml \
  --port $API_GATEWAY_SQL_PORT \
  --enable-https $API_GATEWAY_SQL_ENABLE_HTTPS \
  --cert-file /etc/api-gateway-sql/tls/server.crt \
  --key-file /etc/api-gateway-sql/tls/server.key