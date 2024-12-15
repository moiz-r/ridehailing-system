# Start with the official Golang image
FROM golang:1.23.1-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod, go.sum, and common folder into the build context
COPY user-service/go.mod user-service/go.sum ./

# Download dependencies
RUN go mod download

COPY ../common ./common

# Copy the configs directory
COPY user-service/configs ./configs
# Copy the service's code
COPY . .

# Build the application
RUN go build -o user-service-build user-service/cmd/main.go

# Expose the gRPC port
EXPOSE 50051

EXPOSE 9090

# Run the application
CMD ["./user-service-build"]
