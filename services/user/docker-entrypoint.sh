#!/bin/bash

# Migrate
echo "Migrating"
./app/build/migrate

# Seed permissions
echo "Seeding Permissions"
./app/build/migrate

# Start server
echo "Start Server"
./app/build/server
