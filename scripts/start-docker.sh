#!/bin/sh

docker exec -it app go build -o ./bin/app ./src
docker exec -it app ./bin/app