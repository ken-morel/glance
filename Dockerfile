FROM alpine:3.20
FROM golang:1.21 as builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

COPY go.mod go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o /server

FROM gcr.io/distroless/base-debian11 as final

COPY --from=builder /server /server

ENV PORT 3000
EXPOSE $PORT

WORKDIR /app
COPY build/glance-$TARGETOS-$TARGETARCH${TARGETVARIANT} /app/glance

ENTRYPOINT ["/app/glance"]

