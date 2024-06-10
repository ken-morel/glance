FROM alpine:3.20
FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o /glance

FROM gcr.io/distroless/base-debian11 as final

COPY --from=builder /glance /glance

ENV PORT 3000
EXPOSE $PORT


ENTRYPOINT ["/app/glance"]
