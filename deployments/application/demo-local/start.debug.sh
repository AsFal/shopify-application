
# !/usr/bash/bin

docker-compose \
    -f deployments/application/demo-local/docker-compose.application.yaml \
    -f deployments/deepdetect/demo/docker-compose.deepdetect.yaml \
    -f deployments/elasticsearch/demo/docker-compose.elasticsearch.yaml \
    up --build
