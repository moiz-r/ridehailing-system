# ridehailing-system

## Overview

The ridehailing-system is a microservice-based application designed to manage ride-hailing services. It includes services for user management, booking management, and ride management, all built using Go and PostgreSQL.

## Services

### User Service
- Manages user information.
- Provides endpoints for creating, retrieving, and deleting users.

### Booking Service
- Manages ride bookings.
- Provides endpoints for creating and retrieving bookings.

### Ride Service
- Manages ride information.
- Provides endpoints for creating and retrieving rides.

## Architecture

- **Go**: The primary programming language used for developing the services.
- **PostgreSQL**: The database used for storing user, booking, and ride information.
- **gRPC**: Used for communication between microservices.
- **Prometheus**: Used for monitoring and metrics collection.
- **Docker**: Used for containerizing the services.
- **Docker Compose**: Used for orchestrating multi-container Docker applications.

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Running the Application

1. Clone the repository:
    ```sh
    git clone https://github.com/moiz-r/ridehailing-system.git
    cd ridehailing-system
    ```

2. Build and run the services using Docker Compose:
    ```sh
    docker-compose up --build
    ```

3. The services will be available at the following ports:
    - User Service: `localhost:50051` metrics at `:9090`
    - Booking Service: `localhost:50052` metrics at `:9091`
    - Ride Service: `localhost:50053` metrics at `:9092`

## Configuration

The configuration for each service is managed through configs/config.yml.

