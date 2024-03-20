#!/bin/bash

# Array of microservices
microservices=("msvc-orders" "msvc-payments" "msvc-production")

# Array of subdirectories
subdirs=("endpoint" "service" "transport")

# Loop over each microservice
for msvc in ${microservices[@]}; do
  # Create the microservice directory
  #mkdir $msvc

  # Navigate into the microservice directory
  cd $msvc

  # Create the cmd/server directory
  mkdir -p cmd/server 

  # Navigate into the cmd/server directory
  touch cmd/server/main.go

  # Create a pkg directory
  mkdir pkg

  # Navigate into the pkg directory
  cd pkg

  # Loop over each subdirectory
  for subdir in ${subdirs[@]}; do
    # Create the subdirectory
    mkdir $subdir

    # Navigate into the subdirectory
    cd $subdir

    # Create the .go file
    touch $subdir.go

    # Navigate back to the pkg directory
    cd ..
  done

  # create the folders cmd/api directory for each microservice and add the main.go file

  
  # Navigate back to the parent directory
  cd ../..

done

# remove all cmd/api directories
rm -rf */cmd/api