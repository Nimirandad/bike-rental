#!/bin/bash

set -e

DB_PATH="${DB_PATH:-data/bike_rental.db}"
SCHEMA_PATH="internal/database/schema.sql"
SEED_PATH="internal/database/seed.sql"

echo "Starting database migration..."

if [ "$1" = "fresh" ]; then
    echo "Removing existing database..."
    rm -f "$DB_PATH"
    echo "Database removed"
fi

mkdir -p "$(dirname "$DB_PATH")"

echo "Running schema migration..."
sqlite3 "$DB_PATH" < "$SCHEMA_PATH"
echo "Schema created"

if [ "$1" != "no-seed" ]; then
    echo "Seeding database..."
    sqlite3 "$DB_PATH" < "$SEED_PATH"
    echo "Database seeded"
else
    echo "Skipping seed data"
fi

echo ""
echo "Migration completed successfully!"
echo "Database: $DB_PATH"

if command -v sqlite3 &> /dev/null; then
    echo ""
    echo "Database stats:"
    sqlite3 "$DB_PATH" "SELECT 'Users: ' || COUNT(*) FROM users; SELECT 'Bikes: ' || COUNT(*) FROM bikes; SELECT 'Rentals: ' || COUNT(*) FROM rentals;"
fi
