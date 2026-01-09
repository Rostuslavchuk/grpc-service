# SSO Authentication Service

A robust, production-ready Single Sign-On (SSO) authentication service built with Go and gRPC, designed to provide secure user authentication and authorization capabilities for microservices architectures.

## Overview

This service implements enterprise-grade authentication mechanisms including JWT-based token management, role-based access control (RBAC), and secure session handling. Built with performance and scalability in mind, it serves as the central authentication hub for distributed systems.

## Features

- **Secure Authentication**: Password-based authentication with bcrypt hashing
- **JWT Token Management**: Access and refresh token generation with configurable expiration
- **Role-Based Access Control (RBAC)**: Granular permission management with admin and user roles
- **Database Migrations**: Automated schema management with version control
- **Production Ready**: Structured logging, error handling, and configuration management
- **High Performance**: Optimized SQLite storage with connection pooling
- **Microservice Architecture**: Clean architecture with dependency injection

## Architecture

The service follows clean architecture principles with clear separation of concerns:

```
├── cmd/                    # Application entry points
├── internal/
│   ├── app/               # Application layer
│   ├── config/            # Configuration management
│   ├── domain/            # Domain models and business logic
│   ├── grpc/              # gRPC server implementation
│   ├── jwt/               # JWT token utilities
│   ├── services/          # Business services
│   └── storage/           # Data access layer
├── migrations/            # Database schema migrations
└── config/               # Configuration files
```

## Technology Stack

- **Language**: Go 1.24+
- **RPC Framework**: gRPC
- **Database**: SQLite with migrations
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Build Tool**: Task (Taskfile)
- **Logging**: Structured logging with slog

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Task (Taskfile runner)
- SQLite3

### Installation

```bash
# Clone the repository
git clone https://github.com/Rostuslavchuk/grpc-service.git
cd grpc-service

# Install dependencies
go mod download

# Run database migrations
task migrate

# Start the service
task run
```

### Configuration

The service uses YAML configuration files located in the `config/` directory. Default configuration is provided for local development.

### Running the Service

```bash
# Start the gRPC server (default: localhost:44044)
task run

# Run with custom configuration
GRPC_PORT=50051 task run
```

## API Documentation

The gRPC API is defined in the [sso-protos](https://github.com/Rostuslavchuk/sso-protos) repository. Key endpoints include:

- **Register**: User registration
- **Login**: User authentication
- **Refresh**: Token refresh
- **Validate**: Token validation

## Development

### Running Tests

```bash
# Run all tests
task test

# Run tests with coverage
task test-coverage
```

### Database Operations

```bash
# Run migrations
task migrate

# Create new migration
task migration-create name=migration_name
```

### Code Generation

```bash
# Generate protobuf files (if using local protos)
task proto
```

## Security Considerations

- Passwords are hashed using bcrypt with adaptive cost
- JWT tokens use RS256 signing algorithm
- Configuration supports environment variables for sensitive data
- Database connections use prepared statements to prevent SQL injection

## Performance

- Optimized database queries with proper indexing
- Connection pooling for database operations
- Minimal memory footprint with efficient Go patterns
- gRPC binary protocol for reduced network overhead

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License.
