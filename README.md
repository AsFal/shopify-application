# Shopify Application

## Project Context

You can save your work in a text file, link, etc. and have it ready with you for your technical interview (you can also add it on GitHub). 

TASK: Build an image repository.

You can tackle this challenge using any technology you want. This is an open-ended task.

Please provide brief instructions on how to use your application.

Extra Information: You can, if you wish, use frameworks, libraries and external dependencies to help you get faster to the parts you are interested in building, if this helps you; or start from scratch.

Please focus on what interests you the most. If you need inspiration, here are examples of what you can work on. IF you work on these ideas, we recommend choosing only one or two.

SEARCH function
- from characteristics of the images
- from text
- from an image (search for similar images)
ADD image(s) to the repository
- private or public (permissions)

## Project Description


The following application uses Elastic Search, Deep Detect and Amazon S3 to implement a search API over an image repo. The user can search by iamge, serach by text or search by tags. For a descrption of the API routes see the swagger spec in the api folder.


[Design Process](./doc/design.md)

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
