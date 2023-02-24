# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY redis/*.go ./redis/

RUN GOOS=linux GOARCH=amd64 go build -o /docker-gs-ping

EXPOSE 8080

CMD [ "/docker-gs-ping" ]
