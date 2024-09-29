FROM golang:1.22-alpine AS builder
COPY . .
RUN go mod download
RUN go build -o /glance .
ENTRYPOINT ["/glance"]
