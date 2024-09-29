FROM golang:1.22-alpine AS builder
RUN ls
WORKDIR /app
RUN ls
COPY go.mod go.sum ./
RUN ls
RUN go mod download
RUN ls
COPY . .
RUN ls
RUN go build -o glance .
RUN ls
FROM alpine:latest
RUN ls
RUN ls
WORKDIR /app
RUN ls
COPY --from=builder /app/glance /glance
RUN ls
ENTRYPOINT ["/glance"]
