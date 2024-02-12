# syntax=docker/dockerfile:1

FROM golang:1.22.0-alpine3.19 AS Build
RUN apk add git

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /cache

FROM alpine:latest

WORKDIR /
COPY --from=Build /cache /cache

EXPOSE 8080

ENTRYPOINT [ "/cache" ]
