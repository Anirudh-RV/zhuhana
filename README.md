# Zhuhana

Algorithm Trading Platform

# Run

- To Run the projects run: docker compose -f docker-compose/docker-compose-local.yml up -d --build
- To run without cache:

  > > docker compose -f docker-compose/local-docker-compose.yml build --no-cache
  > > docker compose -f docker-compose/local-docker-compose.yml up

- To run kubernetes cluster:

  > skaffold dev
  > skaffold dev --cleanup=true --no-prune=false
  > skaffold dev --no-prune
  > nohup skaffold dev --cleanup=true --no-prune=false
  > tail -f skaffold.log
  > skaffold delete

- To add the Gateway:

  > helm repo add kong https://charts.konghq.com
  > helm repo update

  > helm install kong kong/kong --namespace kong --create-namespace \
  > --set ingressController.installCRDs=false \
  > --set admin.type=ClusterIP \
  > --set proxy.type=LoadBalancer
