#!/bin/sh

# Wait token
while [ ! -f /output/token.txt ]; do
  sleep 1
done

# Set token as ENV
export TEMP_TOKEN=$(cat /output/token.txt)

echo "TOKEN: $TEMP_TOKEN"

# Run executable
exec ./seed
