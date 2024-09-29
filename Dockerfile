FROM golang:1.22-alpine AS builder
COPY . .
# all files are OK here!
RUN ls
RUN go mod download
RUN go build -o glance .
# COPY --from=builder /glance /glance
ENTRYPOINT ["/glance"]
