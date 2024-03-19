
# Microservices Project

## Delivery Fase 5

The documentation for the Delivery of the Fase 5 of the Tech Challenge [microservices](./DELIVERY.md)

## Overview

This project consists of three microservices: `msvc-orders`, `msvc-payments`, and `msvc-production`. These services are orchestrated using Docker Compose for local development and kubernetes for the production.

## Services

### msvc-orders

This service is responsible for managing orders. It exposes its API on port 8080.

### msvc-payments

This service is responsible for handling payments. It exposes its API on port 8081.

### msvc-production

This service is responsible for managing production. It exposes its API on port 8082.

## How to Run

1. Install Docker and Docker Compose on your machine.
2. Clone this repository.
3. Navigate to the directory containing the `docker-compose.yaml` file.
4. Run the following command to build and start the services:

```bash
# if the project public
docker-compose build 
# if the project is private
docker-compose build  --build-arg GITHUB_ACCESS_TOKEN=your_token
```

Replace `your_token` with your actual GitHub access token if the project is private.

if you want to update the gitsubmodules, you can run the following command:

```bash
git submodule update --init --recursive
# run the custom script
./gitupdate.sh
```

## Troubleshooting

If you encounter any issues while running the services, try the following steps:

1. Check the logs of the services for any error messages:

```bash
docker-compose logs
```

2. Ensure that the ports specified in the `docker-compose.yaml` file are not being used by other processes.

3. Ensure that you have sufficient resources (CPU, Memory, Disk space) available on your machine.

4. If you made changes to the Dockerfile or the source code, rebuild the Docker images:

```bash
docker-compose build --no-cache
```

5. If all else fails, you can reset Docker to its default state:

```bash
docker system prune -a
```

This command will remove all Docker objects from your system. Use it with caution.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)