## run.ps1

# Define a function to print usage information
function usage {
    Write-Host "Usage: .\run.ps1 [start-db|destroy-db|start-app|build-app|start-all|destroy-all|logs-all]"
    exit 1
  }
  
  # Check if docker-compose is installed
  if (!(Get-Command docker-compose -ErrorAction SilentlyContinue)) {
    Write-Host "docker-compose is not installed. Aborting."
    exit 1
  }
  
  # Parse command line arguments
  switch ($args[0]) {
    "start-db" {
      # Start the database container
      docker-compose -f ../devsecops/local/docker-compose.yml up -d db
      Start-Sleep -Seconds 5
      docker-compose -f ../devsecops/local/docker-compose.yml logs db
    }
    "destroy-db" {
      # Destroy the database container and its volumes
      docker-compose -f ../devsecops/local/docker-compose.yml down -v db
    }
    "start-app" {
      # Start the application container
      docker-compose -f ../devsecops/local/docker-compose.yml up -d app
      Start-Sleep -Seconds 5
      docker-compose -f ../devsecops/local/docker-compose.yml logs app
    }
    "build-app" {
      # Build and start the application container
      docker-compose -f ../devsecops/local/docker-compose.yml up --build -d app
      Start-Sleep -Seconds 5
      docker-compose -f ../devsecops/local/docker-compose.yml logs app
    }
    "start-all" {
      # Start all containers
      docker-compose -f ../devsecops/local/docker-compose.yml up -d
      Start-Sleep -Seconds 5
      docker-compose -f ../devsecops/local/docker-compose.yml logs db
      docker-compose -f ../devsecops/local/docker-compose.yml logs app
    }
    "destroy-all" {
      # Destroy all containers and their volumes
      docker-compose -f ../devsecops/local/docker-compose.yml down -v
    }
    "logs-all" {
      # Show logs for all containers
      docker-compose -f ../devsecops/local/docker-compose.yml logs db
      docker-compose -f ../devsecops/local/docker-compose.yml logs app
    }
    default {
      # Print usage information if an invalid command is provided
      usage
    }
  }