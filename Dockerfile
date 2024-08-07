# Use the official Golang image with a compatible version  
FROM golang:1.22-alpine AS builder  

# Set the working directory  
WORKDIR /app  

# Copy go.mod and go.sum files  
COPY go.mod go.sum ./  
# Download the dependencies  
RUN go mod download  

# Copy the source code  
COPY . .  

# Build the Go application  
RUN go build -o myapp .  

# Final image stage  
FROM alpine:latest  

# Install curl to download files  
RUN apk add --no-cache curl  

# Set the working directory for the final image  
WORKDIR /app  

# Download the file  
RUN curl -o glance.yml https://raw.githubusercontent.com/ken-morel/glance/main/glance.yml 

# Copy the compiled binary from the builder stage   
COPY --from=builder /app/myapp /myapp  

# Command to run the binary  
ENTRYPOINT ["/myapp"]
