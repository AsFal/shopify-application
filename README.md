# Shopify Application
The following application uses Elastic Search, Deep Detect and Amazon S3 to implement a search API over an image repo. The user can search by iamge, serach by text or search by tags.

The deployment host must be manually configured
- Setup circleci user
- Add circleci user to docker group
- Add circle ci user ssh key to github
- Add circle ci ssh key to user auth key
- Do the following for elastic search
Add to /etc/sysctl.conf
```sh
grep vm.max_map_count /etc/sysctl.conf
vm.max_map_count=262144
```
And run the following so the live system is up to data
```sh
sysctl -w vm.max_map_count=262144
```

## Demo Deployment

I've created entrypoints for 2 demo deployments. They both only require an installation of Docker and Docker Compose on the host system.

The first is a fully local deployment that stores images in a local folder that is mounted the api.
```sh
bash deployments/application/demo-local/start.sh
```

The second uses a Amazon S3 Bucket to store the images, and comes with a VueJs frontend to easily use the 3 query types. This deployment requires AWS credentials for upload.
```sh
bash deployments/application/demo-amazon/start.sh
```

## API Endpoints
Both deployments will start a container with a mapped port on localhost:8080. The API has the following 3 endpoints:
- POST /
    - -F 'image=@<imagepath> (only jpg supported)
- POST /_search/_image
    - -F 'image=@<imagepath> (only jpg supported)
- GET /_search
    - Query Param 'text=full text search'
    - Query Param 'tags=[tag1,tag2,tag3]'

