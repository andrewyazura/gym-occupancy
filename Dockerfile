#syntax=docker/dockerfile:1

FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY * ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app.out

FROM debian:bookworm

WORKDIR /
COPY --from=build-stage /app.out /app.out
USER nonroot:nonroot

ENTRYPOINT ["/app.out"]
