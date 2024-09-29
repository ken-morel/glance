FROM golang:1.22-alpine AS builder
RUN ls
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o glance .
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/glance /glance
ENTRYPOINT ["/glance"]
