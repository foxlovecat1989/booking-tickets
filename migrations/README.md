# Database Migrations

This directory contains database migrations for the tickets application. The migration system provides versioned database schema changes with both up and down migration capabilities.

## Migration File Format

Migrations follow the naming convention: `{version}_{name}.{direction}.sql`

- `version`: Sequential number (001, 002, 003, etc.)
- `name`: Descriptive name of the migration
- `direction`: Either `up` (apply) or `down` (rollback)

### Current Migrations:
- `001_initial_schema.up.sql` - Creates the initial database schema (concerts, sessions, tickets, orders, payments)
- `001_initial_schema.down.sql` - Rolls back the initial schema
- `002_initial_data.up.sql` - Inserts initial test data (concert, session, tickets)
- `002_initial_data.down.sql` - Removes initial test data

## Available Commands

### Using Makefile (Recommended)

```bash
# Apply all pending migrations
make db-migrate

# Connect to database
make db-connect

# Reset database (drop and recreate)
make db-reset

# Full database reset (removes all volumes)
make db-reset-full
```

### Using the Migration Tool directly

```bash
# Build the migration tool
make migrate-build

# Apply all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Show migration status
make migrate-status

# Create a new migration
make migrate-create name=add_user_table

# Reset all migrations
make migrate-reset
```

### Using the CLI directly

```bash
# Build the migration tool
go build -o bin/migrate cmd/migrate/main.go

# Apply all pending migrations
./bin/migrate -command up

# Rollback last N migrations
./bin/migrate -command down -steps 2

# Show migration status
./bin/migrate -command status

# Create a new migration
./bin/migrate -command create -name add_user_table
```

## Creating New Migrations

### Using Makefile (Recommended):
```bash
make migrate-create name=add_user_table
```

### Using the CLI tool:
```bash
make migrate-build
./bin/migrate -command create -name add_user_table
```

This will create two files:
- `migrations/003_add_user_table.up.sql`
- `migrations/003_add_user_table.down.sql`

### Manual creation:
1. Create two files with the naming convention
2. Add your SQL in the `.up.sql` file
3. Add rollback SQL in the `.down.sql` file

## Migration Best Practices

1. **Always include both up and down migrations**
2. **Use transactions in your SQL when possible**
3. **Test both up and down migrations**
4. **Keep migrations small and focused**
5. **Use descriptive names**
6. **Never modify existing migration files that have been applied**
7. **Test migrations in development before applying to production**

## Current Schema

The current database schema includes:

### Tables
- **concerts**: Concert information (id, name, location, description, created_at)
- **concert_sessions**: Concert sessions (id, concert_id, start_time, end_time, venue, number_of_seats, price)
- **tickets**: Individual tickets (id, session_id, status)
- **orders**: Order records (id, status, total_price, created_at)
- **schema_migrations**: Migration tracking (version, dirty, applied_at)

**Note**: The `order_items` and `payments` tables were removed as they are not used in the current application.

### Indexes
- `idx_concert_sessions_concert_id` - Session by concert lookup (foreign key)
- `idx_tickets_session_id` - Tickets by session lookup (used in GetAvailableTicketsBySessionID)
- `idx_tickets_status` - Tickets by status lookup (used in GetAvailableTicketsBySessionID)

**Note**: Only indexes that are actually used by queries are created. Unused indexes have been removed for better performance.

## Example Migration

### Up Migration (`001_initial_schema.up.sql`):
```sql
-- Create concerts table
CREATE TABLE IF NOT EXISTS concerts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    description TEXT,
    created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000
);

-- Create index on concert name
CREATE INDEX IF NOT EXISTS idx_concerts_name ON concerts(name);
```

### Down Migration (`001_initial_schema.down.sql`):
```sql
-- Drop concerts table
DROP TABLE IF EXISTS concerts CASCADE;
```

## Migration Testing

The migration system includes comprehensive tests:

```bash
# Run migration tests
make test-migrations

# Run all tests (includes migration tests)
make test
```

## Database Reset Workflows

### Development Reset
```bash
# Quick reset for development
make db-reset
```

### Full Reset (Production-like)
```bash
# Complete reset including volume cleanup
make db-reset-full
```

### Manual Reset
```bash
# Stop services
make docker-stop

# Remove volumes
docker volume rm tickets_postgres_data

# Start fresh
make docker-run
```

## Troubleshooting

### Common Issues

1. **Migration already applied**: Use `make migrate-status` to check current state
2. **Database connection issues**: Ensure PostgreSQL is running with `make docker-logs`
3. **Permission errors**: Check Docker volume permissions
4. **Migration conflicts**: Use `make db-reset` to start fresh

### Debug Commands

```bash
# Check migration status
make migrate-status

# View database logs
make docker-logs

# Connect to database for manual inspection
make db-connect

# Check Docker container status
docker ps
```

## Integration with Development Workflow

### Typical Development Session

1. **Start fresh database**:
   ```bash
   make db-reset
   ```

2. **Create new migration** (if needed):
   ```bash
   make migrate-create name=add_new_feature
   ```

3. **Apply migrations**:
   ```bash
   make db-migrate
   ```

4. **Test changes**:
   ```bash
   make test
   ```

5. **Reset for next iteration**:
   ```bash
   make db-reset
   ```

This migration system provides a robust foundation for database schema management in the tickets application. 