#!/bin/bash

# Prompt the user for input
echo -n "Do you want to kill the service? (Y/n) "
read input

# Convert the input to uppercase
input=$(echo "$input" | tr '[:lower:]' '[:upper:]')

# Check if the input is "Y" or "N"
if [ "$input" = "Y" ]; then
    docker kill $(docker ps -q)
elif [ "$input" = "N" ]; then
    echo "Continuing..."
else
    echo "Invalid input. Please enter Y or N."
    exit;
fi

sudo certbot certonly --standalone -d lachiegrant.xyz -d www.lachiegrant.xyz

docker-compose up -d
