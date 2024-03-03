#!/bin/bash
export $(grep -v '^#' scripts/.env | xargs -d '\n')
./scripts/deps-up.sh
go run cmd/main.go