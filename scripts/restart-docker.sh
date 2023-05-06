#!/bin/bash

docker stop app
docker start app
docker exec app ./scripts/start.sh