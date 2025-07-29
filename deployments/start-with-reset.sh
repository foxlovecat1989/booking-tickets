#!/bin/bash

# Start application with automatic database reset
# This script resets the database to initial state before starting the application

echo "Starting application with database reset..."

# Stop any running containers
cd deployments && docker-compose down

# Remove the database volume to ensure clean state
docker volume rm tickets_postgres_data 2>/dev/null || true

# Start PostgreSQL and wait for it to be ready
cd deployments && docker-compose up -d postgres

echo "Waiting for PostgreSQL to be ready..."
until docker-compose exec -T postgres pg_isready -U postgres -d tickets_db; do
  sleep 2
done

echo "PostgreSQL is ready. Starting application..."

# Start the API application
cd deployments && docker-compose up -d api

echo "Application started successfully!"
echo "Database has been reset to initial state with:"
echo "- 1 concert (Rock Concert 2024)"
echo "- 1 concert session (Dec 31, 2024)"
echo "- 10,000 available tickets"
echo ""
echo "API is running on http://localhost:8080"
echo "Database is accessible on localhost:5432"