# Gin Clean Architecture

A clean architecture boilerplate for Golang REST API development using Gin framework. This project follows clean architecture principles to ensure maintainable, testable, and scalable code structure.


## 📋 Features

- **Clean Architecture**: Organized in layers with clear separation of concerns
- **Go with Gin Framework**: High-performance RESTful API with middleware support
- **PostgreSQL Database**: Reliable relational database storage
- **GORM ORM**: Simplified database operations and model mapping
- **Swagger Documentation**: Auto-generated API documentation
- **Database Migrations**: Versioned database schema changes
- **Docker Support**: Containerized development and deployment
- **Hot Reload**: Development with automatic code reloading
- **JWT Authentication**: Secure API endpoints with JWT tokens
- **Environment Configuration**: Easy configuration using .env files

## 🚀 Getting Started

### Prerequisites

- Go 1.16+ installed
- PostgreSQL database
- Docker and Docker Compose (optional, for containerized setup)
- [Golang-migrate](https://github.com/golang-migrate/migrate) (for database migrations)
- [Air](https://github.com/cosmtrek/air) (optional, for hot reloading)
- [Swag](https://github.com/swaggo/swag) (for API documentation)

### Installation

1. Clone the repository
   ```bash
   git clone https://github.com/HasanNugroho/gin-clean.git
   cd gin-clean
   ```

2. Set up the project
   ```bash
   make setup
   ```
   This will download dependencies and create a `.env` file from the example.

3. Configure your `.env` file with your PostgreSQL connection details:
   ```
   POSTGRES_HOST=localhost
   POSTGRES_USER=postgres
   POSTGRES_PASS=postgres
   POSTGRES_DB=gin_clean
   POSTGRES_PORT=5432
   POSTGRES_SSL=disable
   
   JWT_SECRET=your_jwt_secret_key
   ```

4. Set up the database and run migrations
   ```bash
   make docker-run      # Start PostgreSQL container
   make migration-up    # Apply migrations
   ```

5. Build and run the application
   ```bash
   make build
   make run
   ```

## 🛠️ Development

### Hot Reloading

For development with automatic reloading:

```bash
make watch
```

This requires [Air](https://github.com/cosmtrek/air) to be installed.

### Database Migrations

Create a new migration:
```bash
make migration-create desc=add_users_table
```

Apply migrations:
```bash
make migration-up
```

Rollback migrations:
```bash
make migration-down
```

### API Documentation

Generate Swagger documentation:
```bash
make gen-docs
```

View the documentation by running the server and visiting `/swagger/index.html`.

### Authentication

The API uses JWT for authentication. To get a token:

1. Register a new user using the `/api/v1/users/register` endpoint
2. Login with the user credentials at `/api/v1/users/login` to receive a JWT token
3. Use this token in the Authorization header for protected endpoints:
   ```
   Authorization: Bearer your_jwt_token
   ```

## 🐳 Docker

Start the application with Docker:
```bash
make docker-run
```

Stop Docker containers:
```bash
make docker-down
```

## 📚 Project Structure

The project follows clean architecture principles with the following directory structure:

```
.
├── cmd                          # App entry point, runs API server
│   └── api
├── config                       # App configuration
├── container                    # Dependency injection setup
├── docs                         # Documentation and API specs
├── internal                     # Core logic and implementation
│   ├── domain                   # Business logic & entities
│   │   ├── entity               # Domain entities
│   │   ├── repository           # Data access interfaces
│   │   └── service              # Business services
│   ├── infrastructure           # External tools & DB impl.
│   │   └── presistence          # Data persistence layer
│   ├── interfaces               # Communication interfaces
│   │   └── http
│   │       ├── dto              # API DTOs
│   │       ├── handler          # HTTP handlers
│   │       └── middleware       # HTTP middleware
│   └── service                  # Internal helpers & utilities
├── migrations                   # DB migration scripts
└── pkg                          # Shared libraries & utilities

```

### Clean Architecture Layers

- **Entities (Domain)**: Enterprise business rules and objects
- **Use Cases (Service)**: Application-specific business rules
- **Interface Adapters (Delivery)**: Adapters that convert data for use cases and entities
- **Frameworks and Drivers (Infrastructure)**: External frameworks, tools, and delivery mechanisms

## 🧪 Testing

Run unit tests:
```bash
go test ./...
```

## 📝 Copyright
Copyright (c) 2025 Burhan Nurhasan Nugroho.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Swag](https://github.com/swaggo/swag)
- [Air](https://github.com/cosmtrek/air)
