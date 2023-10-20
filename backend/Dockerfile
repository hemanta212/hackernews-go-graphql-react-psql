FROM golang:1.19-alpine as dev
RUN apk add --no-cache make\
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
WORKDIR /go/src/backend
ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

FROM dev as builder
COPY . ./
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/backend/build/hackernews-go-graphql ./
COPY --from=builder /go/src/backend/internal/pkg/db/migrations/postgresql ./internal/pkg/db/migrations/postgresql
CMD ["./hackernews-go-graphql"]
