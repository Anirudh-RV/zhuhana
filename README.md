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
