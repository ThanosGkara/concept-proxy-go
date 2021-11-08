FROM golang:1.17.2-alpine3.14 as builder

# installing git
RUN apk update && apk upgrade && \
    apk add --no-cache git

# setting working directory
WORKDIR /go/src/app

# installing dependencies
RUN go mod init
RUN go get gopkg.in/yaml.v2
RUN go get github.com/patrickmn/go-cache

COPY / /go/src/app/

RUN go build -o /go/src/app/go-proxy main.go

FROM alpine:3.14

ARG MY_SERVICE_PORT=8080

WORKDIR /opt/go-proxy/
COPY --from=builder /go/src/app/go-proxy /opt/go-proxy/go-proxy

EXPOSE ${MY_SERVICE_PORT}

CMD ["/opt/go-proxy/go-proxy"]