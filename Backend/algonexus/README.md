# algonexus

- To Run the projects run (local env):

  > `docker compose --env-file env/local-env.env -f docker-compose/docker-compose-local.yml up -d --build`
  >
  > `docker compose --env-file env/local-env.env -f docker-compose/docker-compose-local.yml up --build`

- To Run the projects run (prod env):

  > `docker compose --env-file env/prod-env.env -f docker-compose/docker-compose-local.yml up -d --build`
  > 
  > `docker compose --env-file env/prod-env.env -f docker-compose/docker-compose-local.yml up --build`

- To Stop:

  > docker compose --env-file env/local-env.env -f docker-compose/docker-compose-local.yml down

- To update Swagger:

  > swag init

- To add spy data, use the container name that is running:
  > kubectl exec -it algonexus-59f75448f4-l85d2 bash
  > run load_spy_data.py
