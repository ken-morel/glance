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
# Install curl to download glance.yml from repository
RUN apk add --no-cache curl
RUN ls
WORKDIR /app
RUN ls
# update the username and repository names to yours please
RUN curl https://raw.githubusercontent.com/ken-morel/glance/main/glance.yml -o glance.yml
RUN ls
COPY --from=builder /app/glance /glance
RUN ls
ENTRYPOINT ["/glance"]
