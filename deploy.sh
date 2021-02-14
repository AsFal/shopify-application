#!/bin/bash

cd ~/shopify-application

git pull

docker kill $(docker -q)
docker rm $(docker -aq)

bash deployments/application/demo-local/start.sh
