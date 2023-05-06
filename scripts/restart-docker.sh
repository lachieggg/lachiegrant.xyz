#!/bin/sh

docker stop webserver
docker cp ./public etc/nginx/html/public
docker start webserver
docker-compose up --build nginx -d
docker stop app
docker start app
docker exec app ./scripts/start.sh