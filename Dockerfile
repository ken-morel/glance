# Use the official Golang image with a compatible version  
FROM golang:1.22-alpine AS builder  

# Set the working directory
COPY glance.yml ./ || COPY ./glance.yml ./ || COPY /glance.yml ./ 
WORKDIR /app  

# Copy go.mod and go.sum files  
COPY go.mod go.sum ./  
# Download the dependencies  
RUN go mod download  
RUN dir -s
# Copy the source code  
COPY . .  
# Build the Go application
RUN go build -o myapp .
RUN dir -s
# Final image stage  
FROM alpine:latest  

# Copy the compiled binary from the builder stage 
#COPY glance.yml /myapp
COPY --from=builder /app/myapp /myapp  
WORKDIR /
# Command to run the binary  
ENTRYPOINT ["/myapp"]
