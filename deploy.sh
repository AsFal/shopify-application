#!/bin/bash

cd ~/shopify-application

git pull


docker kill $(docker ps -q)

docker system prune --force
# Temporary solution because elastic search crashes if volumes are not empty
# This also resets the data 
docker volume prune --force

source ../.aws_cred && bash deployments/application/demo-amazon/start.sh
