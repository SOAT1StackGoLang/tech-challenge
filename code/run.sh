#!/bin/bash

# Define a function to print usage information
usage() {
  echo "Usage: $0 [start-db|destroy-db|start-app|build-app|start-all|destroy-all|recreate-all|recreate-all-with-tests|recreate-all-with-tests-no-build|logs-all|logs-tail]" 1>&2
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

function tag_docker_image() {
  if [ $# -ne 2 ]; then
    echo "Usage: tag_docker_image <origin> <destination>"
    return 1
  fi

  origin=$1
  destination=$2

  if command -v podman >/dev/null 2>&1; then
    tag_cmd="podman"
  elif command -v docker >/dev/null 2>&1; then
    tag_cmd="docker"
  else
    echo "Neither docker nor podman is installed. Aborting."
    exit 1
  fi

  echo "Current images"
  $tag_cmd images
  echo "Tagging $origin as $destination"
  $tag_cmd pull $origin
  $tag_cmd tag $origin $destination
  echo "New images"
  $tag_cmd images
}

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
    sleep 5
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
  recreate-all-with-tests)
    # Destroy all containers and their volumes
    $compose_cmd -f ../devsecops/local/docker-compose.yml down -v
    # Start all containers
    $compose_cmd -f ../devsecops/local/docker-compose.yml up -d --build
    sleep 5
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs db &
    sleep 2
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs app &
    sleep 10
    ./autotest.sh
    ;;
  recreate-all-with-tests-no-build)
    # Destroy all containers and their volumes
    tag_docker_image "local_app:latest" "local-app-image-retag:latest" 
    #tag_docker_image "ghcr.io/soat1stackgolang/tech-challenge:debug-develop" "local-app-image-retag:latest" 

    $compose_cmd -f ../devsecops/local/docker-compose.yml down -v
    # Start all containers
    $compose_cmd -f ../devsecops/local/docker-compose.yml up --no-build -d db
    sleep 5
    $compose_cmd -f ../devsecops/local/docker-compose.yml up --no-build -d app
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs db &
    sleep 2
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs app &
    sleep 10
    ./autotest.sh
    ;;
  recreate-all-with-build-tests-cicd)
    # Build and test image on cicd, uses the local build process, some behaivour may be different 
    # will be builded everytime because of the wierd compose behaviour


    $compose_cmd -f ../devsecops/local/docker-compose.yml down -v
    # Start all containers
    $compose_cmd -f ../devsecops/local/docker-compose.yml up -d db
    sleep 5
    $compose_cmd -f ../devsecops/local/docker-compose.yml up -d app
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs db &
    sleep 2
    $compose_cmd -f ../devsecops/local/docker-compose.yml logs app &
    sleep 10
    ./autotest.sh
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