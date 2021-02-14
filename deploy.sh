#/bin/bash

cd ~/shopify-application

docker kill ($docker -q)
docker rm ($docker -aq)

bash deployments/application/demo-local/start.sh
