#!/bin/sh

current_dir="$(basename "$PWD")"

# Check if we are in the root project folder
if [ "$current_dir" = "scripts" ] || [ ! -f "Makefile" ]; then
    echo "Must be run from the root project folder."
    exit 1
fi

# Prompt the user for input
echo -n "Do you want to kill the service? (y/n) "
read input

# Convert the input to uppercase
input=$(echo "$input" | tr '[:lower:]' '[:upper:]')

# Check if the input is "Y" or "N"
if [ "$input" = "Y" ]; then
    docker kill $(docker ps -q) 2>/dev/null || true
elif [ "$input" = "N" ]; then
    echo "Continuing..."
else
    echo "Invalid input. Please enter Y or N."
    exit 1
fi

# Backup
echo "Backing up..."
sudo mkdir -p /etc/letsencrypt/live/backup
sudo rsync -av --exclude='backup/' /etc/letsencrypt/live/ /etc/letsencrypt/live/backup/

# Delete old certs (but keep the backup directory)
echo "Deleting old certs..."
sudo find /etc/letsencrypt/live/ -mindepth 1 -maxdepth 1 ! -name 'backup' -exec rm -rf {} +

# Update certificate
echo "Updating certificates..."
sudo certbot certonly --standalone -d lachiegrant.xyz -d www.lachiegrant.xyz

# Create directories if they don't exist
echo "Creating certificate directories..."
sudo mkdir -p ./tls/nginx/certs
sudo mkdir -p ./tls/nginx/keys

# Copy new certificate to web server configuration folder
# Find the latest certificate directory
CERT_DIR=$(sudo find /etc/letsencrypt/live -maxdepth 1 -type d -name 'lachiegrant.xyz*' | sort -V | tail -1)

if [ -z "$CERT_DIR" ]; then
    echo "Error: Certificate directory not found!"
    exit 1
fi

echo "Found certificate at: $CERT_DIR"

# Create directories if they don't exist
echo "Creating certificate directories..."
sudo mkdir -p ./tls/nginx/certs
sudo mkdir -p ./tls/nginx/keys

# Copy new certificate to web server configuration folder
echo "Copying new certificate to web server..."
sudo cp "$CERT_DIR/fullchain.pem" ./tls/nginx/certs/fullchain.pem
sudo cp "$CERT_DIR/privkey.pem" ./tls/nginx/keys/privkey.pem
sudo chown -R $(whoami) ./tls/nginx/

# Restart service
echo "Restarting docker-compose..."
docker-compose up -d
