#!/bin/bash

set -e

migrate -database "$DATABASE_URL" -path /app/database/migrations up

