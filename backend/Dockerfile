FROM golang:1.18-alpine

RUN apk add --no-cache make

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

ENV APP_HOME /go/src/backend
RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"
EXPOSE 8080

RUN mkdir -p /etc/letsencrypt/live/vps.osac.org.np

CMD ["make", "run"]
