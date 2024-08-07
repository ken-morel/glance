# First stage: build the Go application  
FROM golang:1.20-alpine AS builder  

ARG TARGETOS  
ARG TARGETARCH  
ARG TARGETVARIANT  

WORKDIR /app  

# Copy the Go module files  
COPY go.mod go.sum ./  

# Download dependencies  
RUN go mod download  

# Copy the source code  
COPY . .  

# Build the "glance" application  
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOARM=$TARGETVARIANT go build -o /glance .  

# Final stage: create the runtime image  
FROM alpine:3.20  

WORKDIR /app  

# Copy the built application from the builder stage  
COPY --from=builder /glance /glance  

EXPOSE 8080/tcp  
ENTRYPOINT ["/glance"]
