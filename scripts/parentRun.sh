#!/bin/bash

######################################
# This run script assumes it is cloned as a gitsubmodule parented by
# the Chat project.  This will run this project on localhost but will expect
# the parent dependencies and services to be running with the source env vars.
######################################

export $(grep -v '^#' ../.env | xargs -d '\n')
go run cmd/main.go