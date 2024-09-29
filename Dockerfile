FROM golang:1.22-alpine AS builder
RUN ls 
# first
RUN go mod download
RUN go build -o glance .
FROM alpine:latest
COPY --from=builder /glance /glance
ENTRYPOINT ["/glance"]
