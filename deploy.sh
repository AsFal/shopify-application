#!/bin/bash

cd ~/shopify-application

git pull

docker kill $(docker ps -q)
docker rm $(docker ps -aq)

bash deployments/application/demo-local/start.sh
