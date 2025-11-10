#!/bin/zsh
export GOOSE_DRIVER=sqlite3
export GOOSE_DBSTRING=./data/database.sqlite
export GOOSE_MIGRATION_DIR=./data/migrations
echo "loading env vars..."
