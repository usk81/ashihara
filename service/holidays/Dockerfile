# Build
FROM golang:1.22-alpine3.19 AS build

WORKDIR /go/app

COPY go.mod go.sum ./
RUN apk --no-cache add git && \
  go mod download
COPY ./ ./
RUN set -eux && \
  go build -o api main.go

# Production
FROM alpine:3.19

RUN apk --no-cache add tzdata mysql-client && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENV TZ Asia/Tokyo

WORKDIR /app

COPY --from=build /go/app/api .

RUN ls -la 

ENTRYPOINT ["/app/api"]