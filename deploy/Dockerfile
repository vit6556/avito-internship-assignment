FROM golang:1.23.2

WORKDIR ${GOPATH}/avito-shop/
COPY . ${GOPATH}/avito-shop/

RUN go build -o /http-server ./cmd/http-server \
    && go clean -cache -modcache

RUN go build -o /migrator ./cmd/migrator \
    && go clean -cache -modcache

COPY deploy/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]