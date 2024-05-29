# syntax=docker/dockerfile:1

FROM golang:1.22-alpine3.19 AS Build
RUN apk add git

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /customers

FROM alpine:latest

WORKDIR /
COPY --from=Build /customers /customers

EXPOSE 8080

ENTRYPOINT [ "/customers" ]
