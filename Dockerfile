FROM golang:1.21 AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o api-gateway-sql .

FROM alpine:3.20

COPY --from=builder /build/api-gateway-sql /usr/local/bin
COPY --from=builder /build/entrypoint.sh /usr/local/bin
RUN chmod +x /usr/local/bin/entrypoint.sh

RUN mkdir /data
RUN mkdir /etc/api-gateway-sql
RUN mkdir /etc/api-gateway-sql/tls
COPY --from=builder --chown=nobody /build/fixtures/config.default.yaml /etc/api-gateway-sql/config.yaml
COPY --from=builder --chown=nobody /build/fixtures/tls/server.crt /etc/api-gateway-sql/tls/server.crt
COPY --from=builder --chown=nobody /build/fixtures/tls/server.key /etc/api-gateway-sql/tls/server.key

RUN apk update && apk add --no-cache ca-certificates

ENV API_GATEWAY_SQL_PORT=5297
ENV API_GATEWAY_SQL_ENABLE_HTTPS="true"

USER nobody
EXPOSE $API_GATEWAY_SQL_PORT

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]