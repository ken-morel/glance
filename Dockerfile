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

# Final stage  
FROM alpine:3.20  

ARG TARGETOS  
ARG TARGETARCH  
ARG TARGETVARIANT  

WORKDIR /app  
COPY --from=builder /glance /glance  

EXPOSE 8080/tcp  
ENTRYPOINT ["/glance"]
