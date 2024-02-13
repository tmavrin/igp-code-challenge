#!/bin/bash

command="$1"

if [[ $command == "debug" ]]; then
    docker compose -f docker-compose-debug.yaml up --build
else 
    echo "Unknown command: $command"
fi