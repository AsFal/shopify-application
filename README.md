# Shopify Application
## Demo Deployment
```bash
docker-compose $(find deployments/demo/docker-compose* | sed -e 's/^/-f /') up -d
```