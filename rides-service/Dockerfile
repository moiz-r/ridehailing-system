# Start with the official Golang image
FROM golang:1.23.1-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod, go.sum, and common folder into the build context
COPY rides-service/go.mod rides-service/go.sum ./

# Download dependencies
RUN go mod download

COPY ../common ./common

# Copy the configs directory
COPY rides-service/configs ./configs
# Copy the service's code
COPY . .

# Build the application
RUN go build -o rides-service-build rides-service/cmd/main.go

# Expose the gRPC port
EXPOSE 50053

EXPOSE 9092

# Run the application
CMD ["./rides-service-build"]
