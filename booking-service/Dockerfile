# Start with the official Golang image
FROM golang:1.23.1-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod, go.sum, and common folder into the build context
COPY booking-service/go.mod booking-service/go.sum ./

# Download dependencies
RUN go mod download

COPY ../common ./common

# Copy the configs directory
COPY booking-service/configs ./configs
# Copy the service's code
COPY . .

# Build the application
RUN go build -o booking-service-build booking-service/cmd/main.go

# Expose the gRPC port
EXPOSE 50052

EXPOSE 9091

# Run the application
CMD ["./booking-service-build"]
