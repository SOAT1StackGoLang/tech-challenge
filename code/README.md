# Code

This directory contains the application code and tests.

## Directory Structure

```bash
code/
├── cmd
│   ├── migrations
│   │   └── files
│   └── web
├── helpers
└── internal
    ├── core
    │   ├── domain
    │   ├── ports
    │   └── usecases
    ├── handlers
    │   └── http
    ├── repositories
    │   └── postgres
    └── setup
        └── apidocs
```

# Scripts

## run.sh

The `run.sh` script is a Bash script that provides a simple interface to manage the containers defined in the `docker-compose.yml` file.

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
- `recreate-all`: Destroys all containers and their volumes and starts all containers.
- `recreate-all-with-tests`: Destroys all containers and their volumes, starts all containers, and runs the autotest script.
- `recreate-all-with-tests-no-build`: Destroys all containers and their volumes, then starts all containers and runs the autotest script without building the application container, pulling the image from the current commit-sha(must being a builded commit).
- `logs-all`: Shows the logs of all containers.
- `logs-tail`: Shows the logs of all containers and follows the output.

For example, to start the database container, run:

```shellscript
./run.sh start-db
```

This will start the database container and show its logs.

To recreate all containers and run the autotest script, run:

```shellscript
./run.sh recreate-all-with-tests
```

or if your commit-sha was previous builded you can run

```shellscript
./run.sh recreate-all-with-tests-no-build
```

files locations for the scripts:

```bash
devsecops/local/
├── README.md
├── code
│   └── Dockerfile
├── docker-compose.yml
└── initdb
    └── 00-configs.sql
```

## autotest.sh

[Autotest](AUTOTEST.md)
