#!/bin/sh

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h postgres -p 5432 -U postgres; do
  sleep 1
done

echo "PostgreSQL is ready. Running migrations..."

# Run migrations
for file in /app/migrations/*.up.sql; do
  echo "Running migration: $file"
  psql -h postgres -p 5432 -U postgres -d web3_edu_db -f "$file"
done

echo "Migrations completed."

# Start the application
exec "$@"
