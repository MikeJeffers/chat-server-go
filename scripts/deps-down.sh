#!/bin/bash
export $(grep -v '^#' scripts/.env | xargs -d '\n')
docker compose down -v
