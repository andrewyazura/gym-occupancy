#syntax=docker/dockerfile:1

FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
ADD config ./config
ADD gymportal ./gymportal
ADD storage ./storage

RUN CGO_ENABLED=0 GOOS=linux go build -o /app.out

FROM debian:bookworm
RUN groupadd -r deploy && useradd -r -g deploy deploy

RUN apt update && \
    apt upgrade -y && \
    apt install -y ca-certificates

USER deploy

WORKDIR /home/deploy
COPY --from=build-stage /app.out /app.out

ENTRYPOINT ["/app.out"]
