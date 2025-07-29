# Tickets API

A gRPC-based ticket management system built with Go, following standard Go project layout conventions. This system provides a complete solution for managing concert tickets, orders, and availability tracking.

## 🚀 Features

- **gRPC API**: High-performance RPC communication with Protocol Buffers
- **PostgreSQL**: Reliable data storage with ACID compliance
- **Domain-Driven Design**: Clean architecture with clear separation of concerns
- **Docker Support**: Containerized deployment with Docker Compose
- **Database Migrations**: Version-controlled schema changes
- **Comprehensive Testing**: Unit, integration, and gRPC handler tests
- **Structured Logging**: Configurable logging with different levels and formats
- **Transaction Management**: Robust database transaction handling
- **Protocol Buffer Development**: Automated code generation with cleanup

## 🏗️ Project Structure

```
tickets/
├── api/                   # Generated API code (protobuf)
│   ├── tickets.pb.go     # Generated protobuf messages
│   └── tickets_grpc.pb.go # Generated gRPC service code
├── proto/                 # Protocol Buffer definitions
│   └── tickets.proto     # gRPC service definitions
├── cmd/                   # Main applications
│   ├── server/           # gRPC server application (currently setup only)
│   └── migrate/          # Database migration tool
├── internal/              # Private application and library code
│   ├── config/           # Configuration management
│   ├── handler/          # gRPC request/response handlers
│   ├── logger/           # Structured logging
│   ├── migrations/       # Migration management
│   ├── models/           # Domain models
│   │   ├── domain/       # Business domain models
│   │   └── db/           # Database models
│   ├── repository/       # Data access layer
│   └── service/          # Business logic layer
├── migrations/            # Database migration files
├── deployments/           # Deployment configurations
├── config.yaml           # Application configuration
├── Makefile              # Build and development commands
└── go.mod                # Go module dependencies
```

## 🎯 Implemented Services

### ✅ Currently Available
- **Order Service**: Create ticket orders with transaction safety
- **Concert Session Repository**: Manage concert sessions and availability
- **Ticket Repository**: Handle ticket operations and status tracking
- **gRPC Handler**: Complete gRPC endpoint implementation with error handling
- **Database Migrations**: Version-controlled schema management
- **Structured Logging**: Configurable logging with Logrus
- **Configuration Management**: Environment-based configuration
- **Server Setup**: Database connection and migration initialization

### 🔄 Planned Services
- **gRPC Server**: ✅ Server now starts and listens on configured port
- **GetOrder**: Retrieve order details by ID
- **ListOrders**: List user orders with pagination
- **GetConcertSession**: Get concert session details
- **ListConcertSessions**: List available sessions
- **GetAvailableTickets**: Get available tickets for a session
- **Payment Service**: Handle payment processing
- **User Service**: User management and authentication
- **Health Service**: Service health monitoring

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Docker & Docker Compose (optional)
- Protocol Buffer compiler (`protoc`)

### Local Development

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd tickets
   make deps
   ```

2. **Generate Protocol Buffer code**:
   ```bash
   make proto
   ```

3. **Start database**:
   ```bash
   make docker-run
   ```

4. **Run migrations**:
   ```bash
   make db-migrate
   ```

5. **Build and run server**:
   ```bash
   make build
   make run
   ```

**Note**: The server now starts the gRPC server on the configured port and is ready to handle requests.

### Docker Deployment

```bash
# Build and run with Docker Compose
make docker-build
make docker-run

# View logs
make docker-logs

# Stop services
make docker-stop
```

## 📡 gRPC API

The system defines a gRPC API with the following endpoints:

### Order Management
- `CreateOrder`: ✅ Handler implemented and server running
- `GetOrder`: Retrieve order details (planned)
- `ListOrders`: List user orders (planned)

### Concert Management
- `GetConcertSession`: Get concert session details (planned)
- `ListConcertSessions`: List available sessions (planned)
- `GetAvailableTickets`: Get available tickets for a session (planned)

### Example gRPC Request
```protobuf
// Create an order
CreateOrderRequest {
  user_id: 1
  concert_session_id: 1
  number_of_tickets: 2
}
```

### Example gRPC Response
```protobuf
CreateOrderResponse {
  order_id: 123
  status: "pending"
  ticket_ids: ["uuid-1", "uuid-2"]
  total_price: 199.98
  created_at: "2024-12-31T20:00:00Z"
}
```

## 🔧 Development

### Protocol Buffer Development

```bash
# Generate Protocol Buffer code (with automatic cleanup)
make proto

# Clean generated files manually
make proto-clean

# View all available commands
make help
```

The `make proto` command automatically:
- Cleans existing generated files
- Generates fresh Go code from `.proto` files
- Organizes files in the proper directory structure
- Provides clear progress feedback

### Code Organization

- **Domain Layer** (`internal/models/domain/`): Core business models and entities
- **Repository Layer** (`internal/repository/`): Data access and persistence
- **Service Layer** (`internal/service/`): Business logic and orchestration
- **Handler Layer** (`internal/handler/`): gRPC request/response handling
- **Configuration** (`internal/config/`): Application configuration
- **API Layer** (`api/`): Generated Protocol Buffer code

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean
```

### Database Operations

```bash
# Run migrations
make db-migrate

# Connect to database
make db-connect

# Reset database
make db-reset
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run specific test categories
make test-models
make test-repository
make test-service
make test-config
make test-logger
make test-migrations
```

### Test Coverage

The project includes comprehensive test coverage:
- **Unit Tests**: 65 tests covering models, config, and logger
- **Integration Tests**: 24 tests covering repository and service layers
- **gRPC Handler Tests**: 12 tests covering API endpoints
- **Total**: 101 tests with detailed coverage reporting

## ⚙️ Configuration

The application uses `config.yaml` for configuration:

```yaml
server:
  port: 8080
  grpc_port: 9090

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "password"
  dbname: "tickets_db"
  url: "postgres://postgres:password@localhost:5432/tickets_db?sslmode=disable"

logging:
  level: "info"
  format: "text"
  output: "stdout"
  include_caller: true
  include_timestamp: true

mode: "debug"
port: "8080"
```

## 🗄️ Database Schema

The system includes the following core tables:

- **concerts**: Concert information (name, location, description)
- **concert_sessions**: Concert sessions with pricing and timing
- **tickets**: Individual tickets with availability status
- **orders**: Order records with status and pricing
- **order_items**: Order-ticket relationships
- **payments**: Payment records and status
- **schema_migrations**: Migration tracking table

## 📚 Documentation

- [Testing Guide](TESTING_GUIDE.md) - Comprehensive testing documentation
- [Logging Guide](LOGGING_GUIDE.md) - Logging configuration and usage
- [Migration Guide](MIGRATION_GUIDE.md) - Database migration documentation
- [Migrations README](migrations/README.md) - Migration system details

## 🛠️ Available Make Commands

### Build & Run
- `make build` - Build the application
- `make run` - Run the application (setup only)
- `make clean` - Clean build artifacts

### Protocol Buffer
- `make proto` - Generate Protocol Buffer code (with cleanup)
- `make proto-clean` - Remove generated Protocol Buffer files

### Testing
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report
- `make test-unit` - Run unit tests only
- `make test-integration` - Run integration tests only

### Database
- `make db-migrate` - Run database migrations
- `make db-connect` - Connect to database
- `make db-reset` - Reset database to initial state

### Docker
- `make docker-build` - Build Docker image
- `make docker-run` - Start services with Docker Compose
- `make docker-stop` - Stop services
- `make docker-logs` - View logs

### Development
- `make deps` - Install dependencies
- `make fmt` - Format code
- `make lint` - Run linter
- `make help` - Show all available commands

## 🤝 Contributing

1. Follow Go conventions and project structure
2. Write tests for new features
3. Update documentation as needed
4. Use conventional commit messages
5. Ensure all tests pass before submitting
6. Run `make proto` after modifying `.proto` files

## 📄 License

This project is licensed under the MIT License. 