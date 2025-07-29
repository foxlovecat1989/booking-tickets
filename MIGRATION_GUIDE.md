# Database Migration System

This document provides a comprehensive guide to the database migration system implemented for the tickets application.

## Overview

The migration system provides:
- **Versioned migrations** with sequential numbering
- **Up and down migrations** for applying and rolling back changes
- **CLI tool** for managing migrations
- **Integration helpers** for application startup
- **Transaction safety** with automatic rollback on failure
- **Migration tracking** in a dedicated `schema_migrations` table

## Quick Start

### 1. Build the migration tool
```bash
make migrate-build
```

### 2. Apply all pending migrations
```bash
make migrate-up
```

### 3. Check migration status
```bash
make migrate-status
```

### 4. Create a new migration
```bash
make migrate-create name=add_user_table
```

## Architecture

### Core Components

1. **MigrationManager** (`internal/migrations/migration.go`)
   - Handles loading, parsing, and executing migrations
   - Manages migration state and tracking
   - Provides transaction safety

2. **CLI Tool** (`cmd/migrate/main.go`)
   - Command-line interface for migration operations
   - Supports up, down, status, and create commands
   - Uses project configuration

3. **Integration Helpers** (`internal/migrations/integration.go`)
   - Simple functions for application integration
   - Automatic migration on startup

### Migration File Structure

```
migrations/
├── 001_initial_schema.up.sql      # Apply initial schema
├── 001_initial_schema.down.sql    # Rollback initial schema
├── 002_add_users.up.sql           # Add users table
├── 002_add_users.down.sql         # Remove users table
└── README.md                      # Migration documentation
```

## Migration File Format

### Naming Convention
```
{version}_{name}.{direction}.sql
```

- `version`: Sequential number (001, 002, 003, etc.)
- `name`: Descriptive name (use underscores for spaces)
- `direction`: `up` (apply) or `down` (rollback)

### File Content Structure

#### Up Migration Example
```sql
-- Migration: add_user_table
-- Version: 2
-- Created: 2024-01-15

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

#### Down Migration Example
```sql
-- Rollback: add_user_table
-- Version: 2
-- Created: 2024-01-15

DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users CASCADE;
```

## Available Commands

### Makefile Commands (Recommended)

| Command | Description |
|---------|-------------|
| `make migrate-build` | Build the migration CLI tool |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Rollback the last migration |
| `make migrate-status` | Show migration status |
| `make migrate-create name=<name>` | Create new migration |
| `make migrate-reset` | Reset all migrations |

### Direct CLI Commands

```bash
# Build tool
go build -o bin/migrate cmd/migrate/main.go

# Apply migrations
./bin/migrate -command up

# Rollback migrations
./bin/migrate -command down -steps 2

# Show status
./bin/migrate -command status

# Create migration
./bin/migrate -command create -name add_payment_methods
```

## Integration with Application

### Automatic Migration on Startup

The server automatically runs migrations on startup:

```go
// In cmd/server/main.go
if err := migrations.RunMigrationsOnStartup(db.DB, "migrations"); err != nil {
    log.Fatalf("Failed to run migrations: %v", err)
}
```

### Manual Integration

```go
import "tickets/internal/migrations"

// Run migrations manually
if err := migrations.RunMigrationsOnStartup(db, "migrations"); err != nil {
    // Handle error
}

// Check migration status
status, err := migrations.GetMigrationStatusOnStartup(db, "migrations")
if err != nil {
    // Handle error
}
```

## Migration Tracking

The system tracks applied migrations in the `schema_migrations` table:

```sql
CREATE TABLE schema_migrations (
    version BIGINT PRIMARY KEY,
    dirty BOOLEAN NOT NULL DEFAULT FALSE,
    applied_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Migration States

- **Applied**: Migration has been successfully applied
- **Pending**: Migration exists but hasn't been applied
- **Dirty**: Migration failed during application (requires manual intervention)

## Best Practices

### 1. Migration Design
- **Keep migrations small and focused**
- **Always include both up and down migrations**
- **Use descriptive names**
- **Test migrations before applying to production**

### 2. SQL Best Practices
- **Use transactions when possible**
- **Handle foreign key constraints properly**
- **Use `IF EXISTS` and `IF NOT EXISTS` clauses**
- **Include proper indexes**

### 3. Version Control
- **Never modify existing migration files that have been applied**
- **Create new migrations for schema changes**
- **Keep migration files in version control**

### 4. Testing
- **Test both up and down migrations**
- **Test with sample data**
- **Verify foreign key relationships**

## Example Workflow

### 1. Create a New Migration
```bash
make migrate-create name=add_user_roles
```

This creates:
- `migrations/002_add_user_roles.up.sql`
- `migrations/002_add_user_roles.down.sql`

### 2. Edit the Migration Files
```sql
-- 002_add_user_roles.up.sql
CREATE TABLE user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
```

```sql
-- 002_add_user_roles.down.sql
DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP TABLE IF EXISTS user_roles CASCADE;
```

### 3. Apply the Migration
```bash
make migrate-up
```

### 4. Verify the Migration
```bash
make migrate-status
```

### 5. Rollback if Needed
```bash
make migrate-down
```

## Troubleshooting

### Common Issues

#### Migration Fails to Apply
- Check SQL syntax
- Verify database connection
- Ensure previous migrations are applied
- Check for foreign key constraints

#### Migration Already Applied
- The system skips already applied migrations
- Check `schema_migrations` table for applied versions

#### Rollback Fails
- Verify down migration SQL is correct
- Check that migration was actually applied
- Ensure foreign key constraints are handled

#### Database Connection Issues
- Verify database is running
- Check connection string in `config.yaml`
- Ensure database user has proper permissions

### Debugging Commands

```bash
# Check migration status
make migrate-status

# View database logs
make docker-logs

# Connect to database
make db-connect

# Reset all migrations (development only)
make migrate-reset
```

## Configuration

The migration system uses the database configuration from `config.yaml`:

```yaml
database:
  url: "postgres://postgres:password@localhost:5432/tickets_db?sslmode=disable"
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "password"
  dbname: "tickets_db"
```

## Testing

Run migration tests:
```bash
go test ./internal/migrations/
```

The tests verify:
- Migration file parsing
- Migration loading
- Status formatting
- Error handling

## Production Considerations

### 1. Backup Strategy
- Always backup database before running migrations
- Test migrations on staging environment first
- Have rollback plan ready

### 2. Zero-Downtime Deployments
- Design migrations to be backward compatible
- Use feature flags for application changes
- Consider blue-green deployment strategy

### 3. Monitoring
- Monitor migration execution time
- Set up alerts for migration failures
- Track migration success rates

### 4. Security
- Use dedicated database user for migrations
- Limit migration user permissions
- Audit migration changes

## Migration Examples

### Adding a New Table
```sql
-- Up migration
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Down migration
DROP TABLE IF EXISTS products CASCADE;
```

### Adding a Column
```sql
-- Up migration
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- Down migration
ALTER TABLE users DROP COLUMN IF EXISTS phone;
```

### Adding an Index
```sql
-- Up migration
CREATE INDEX idx_users_phone ON users(phone);

-- Down migration
DROP INDEX IF EXISTS idx_users_phone;
```

### Data Migration
```sql
-- Up migration
UPDATE users SET status = 'active' WHERE status IS NULL;

-- Down migration
UPDATE users SET status = NULL WHERE status = 'active';
```

This migration system provides a robust, versioned approach to database schema management with full rollback capabilities and integration with your Go application. 