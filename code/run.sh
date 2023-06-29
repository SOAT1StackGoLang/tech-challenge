#!/bin/bash

# Define a function to print usage information
usage() {
  echo "Usage: $0 [start-db|destroy-db|start-app|build-app|start-all|destroy-all|logs-all]" 1>&2
  exit 1
}

# Check if podman-compose is installed, otherwise use docker-compose
if command -v podman-compose >/dev/null 2>&1; then
  compose_cmd="podman-compose"
elif command -v docker-compose >/dev/null 2>&1; then
  compose_cmd="docker-compose"
else
  echo "Neither docker-compose nor podman-compose is installed. Aborting."
  exit 1
fi

# Parse command line arguments
case "$1" in
  start-db)
    # Start the database container
    $compose_cmd -f ../devsecops/local/docker-compose.yml up -d db
    sleep 5
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs db
    ;;
  destroy-db)
    # Destroy the database container and its volumes
    $compose_cmd -f ../devsecops/local/docker-compose.yml down -v db
    ;;
  start-app)
    # Start the application container
    $compose_cmd -f ../devsecops/local/docker-compose.yml up -d app --build
    sleep 5
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs app
    ;;
  build-app)
    # Build and start the application container
    $compose_cmd -f ../devsecops/local/docker-compose.yml up -d app --build
    sleep 5
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs app
    ;;
  start-all)
    # Start all containers
    $compose_cmd -f ../devsecops/local/docker-compose.yml up -d --build
    sleep 5
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs db &
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs app &
    ;;
  destroy-all)
    # Destroy all containers and their volumes
    $compose_cmd -f ../devsecops/local/docker-compose.yml down -v
    ;;
  recreate-all)
    # Destroy all containers and their volumes
    $compose_cmd -f ../devsecops/local/docker-compose.yml down -v
    # Start all containers
    $compose_cmd -f ../devsecops/local/docker-compose.yml up -d --build
    sleep 5
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs db &
    sleep 2
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs app &
    ;;
  logs-all)
    # Show logs for all containers
    if $compose_cmd -f ../devsecops/local/docker-compose.yml ps | grep -q "local_db_1"; then
      $compose_cmd -f ../devsecops/local/docker-compose.yml logs -f db &
    fi
    if $compose_cmd -f ../devsecops/local/docker-compose.yml ps | grep -q "local_app_1"; then
      $compose_cmd -f ../devsecops/local/docker-compose.yml logs -f app &
    fi
    ;;
  logs-tail)
    # Show logs for all containers
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs -f
    ;;
  *)
    # Print usage information if an invalid command is provided
    usage
    ;;
esac

exit 0