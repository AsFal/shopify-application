# !/usr/bash/bin

docker-compose \
    -f deployments/application/demo-amazon/docker-compose.application.yaml \
    -f deployments/deepdetect/default/docker-compose.deepdetect.yaml \
    -f deployments/elasticsearch/demo/docker-compose.elasticsearch.yaml \
    up --build
 