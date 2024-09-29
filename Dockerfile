FROM golang:1.22-alpine AS builder  
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download  
COPY . .  
RUN go build -o myapp .  
FROM alpine:latest
# Install curl to download glance.yml from repository
RUN apk add --no-cache curl  
WORKDIR /app  
# update the username and repository names to yours please
RUN curl https://raw.githubusercontent.com/ken-morel/glance/main/glance.yml -o glance.yml   
COPY --from=builder /app/myapp /myapp  
ENTRYPOINT ["/myapp"]
