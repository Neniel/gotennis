# syntax=docker/dockerfile:1

FROM golang:1.22-alpine3.19 AS Build
RUN apk add git

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /auth

FROM alpine:latest

WORKDIR /
COPY --from=Build /auth /auth

EXPOSE 8080

ENTRYPOINT [ "/auth" ]
