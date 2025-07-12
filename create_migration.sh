#!/bin/bash

migration_name=$1

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <migration_name>"
    exit 1
fi

docker run --rm -v ./migrations:/migrations migrate/migrate:v4.18.3 create \
-ext sql -dir migrations ${migration_name}