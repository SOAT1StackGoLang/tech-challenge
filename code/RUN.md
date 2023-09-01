# First Initiate the containers locally

We can start or recreate from scratch

```sh
## Start the application first time (or after a clean)
./run.sh start-all
```

```sh
## Stop the application and delete all containers, volumes and networks and start again
./run.sh recreate-all
```

```sh
## Stop the application and delete all containers, volumes and networks
./run.sh destroy-all
```

```sh
## Recreate if previous run, otherwise start the application and run tests
## The tests is the script autotest.sh, will test all relevantes endpoints and flows
./run.sh recreate-all-with-tests
```