# Code

This directory contains the application code and tests.


## run.sh

The `run.sh or run.ps1` script is a Bash script that provides a simple interface to manage the containers defined in the `docker-compose.yml` file.

### Dependencies
The script requires either docker-compose or podman-compose to be installed on your system. If neither is installed, the script will not work.

### Usage

To use the script, open a terminal and navigate to the root directory of the project. Then, run the script with one of the following commands:

- `start-db`: Starts the database container.
- `destroy-db`: Destroys the database container and its volumes.
- `start-app`: Starts the application container.
- `build-app`: Builds and starts the application container.
- `start-all`: Starts all containers.
- `destroy-all`: Destroys all containers and their volumes.
- `logs-all`: Shows the logs of all containers.

For example, to start the database container, run:

```shellscript
./run.sh start-db
```

```shellscript
.\run.ps1 start-db
```

This will start the database container and show its logs.