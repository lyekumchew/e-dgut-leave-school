FROM golang:1.15-alpine AS builder

#RUN apk update && apk add build-base gcc

WORKDIR /src

COPY . .

RUN go mod download \
    && GOARCH=amd64 GOOS=linux GOBUILD=CGO_ENABLED=0 go build -ldflags '-w -s' -o edgut-client


FROM alpine:latest

LABEL maintainer="AUTUMN"

WORKDIR /app

COPY --from=builder /src/edgut-client /app/

RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && rm -rf /var/cache/apk/*

ENTRYPOINT ["/app/edgut-client", "-cron"]