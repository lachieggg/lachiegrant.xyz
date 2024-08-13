#!/bin/bash

current_dir="$(basename "$PWD")"

# Check if we are in the root project folder
if [ "$current_dir" = "scripts" ] || [ ! -f "Makefile" ]; then
    echo "Must be run from the root project folder."
    exit;
fi

# Prompt the user for input
echo -n "Do you want to kill the service? (y/n) "
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

# Backup
echo "Backing up..."
sudo mkdir -p /etc/letsencrypt/live/backup
sudo rsync -av --exclude='backup/' /etc/letsencrypt/live/ /etc/letsencrypt/live/backup/

# Delete old certs
echo "Deleting old certs..."
sudo find /etc/letsencrypt/live/ -mindepth 1 -maxdepth 1 ! -name 'backup' -exec rm -rf {} +
sudo find /etc/letsencrypt/live/

# Update certificate
echo "Updating certificates"
sudo certbot certonly --standalone -d lachiegrant.xyz -d www.lachiegrant.xyz

# Move cert to default folder
echo "Move certificate to default folder"
sudo find /etc/letsencrypt/live/ -mindepth 1 -maxdepth 1 -name 'lachiegrant.xyz*' \
	-exec mv {} /etc/letsencrypt/live/lachiegrant.xyz/ \;

# Copy new certificate to web server configuration folder
echo "Copying new certificate to web server"
sudo cp /etc/letsencrypt/live/lachiegrant.xyz/fullchain.pem ./tls/nginx/certs/fullchain.pem
sudo cp /etc/letsencrypt/live/lachiegrant.xyz/privkey.pem ./tls/nginx/keys/privkey.pem
sudo find /etc/letsencrypt/live

# Restart service
docker-compose up -d
