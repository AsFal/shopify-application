#!/bin/bash

cd ~/shopify-application

git pull


docker kill $(docker ps -q)

docker system prune
# Temporary solution because elastic search crashes if volumes are not empty
# This also resets the data 
docker volume prune


bash deployments/application/demo-local/start.sh
