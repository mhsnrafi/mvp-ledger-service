#!/bin/bash
set -e

# Add the services and their respective ports you want to wait for
services="nginx:4000 ledger-service:8000"

for service in $services; do
  host=$(echo $service | cut -d: -f1)
  port=$(echo $service | cut -d: -f2)
  echo "Waiting for $host:$port to become available..."
  while ! nc -z $host $port; do
    sleep 1
  done
  echo "$host:$port is available!"
done

# Run the test command after the services are available
exec "${@:1}"