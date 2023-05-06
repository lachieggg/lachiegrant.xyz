#!/bin/sh

docker-compose up --build nginx -d
docker stop app
docker start app
docker exec app ./scripts/start.sh