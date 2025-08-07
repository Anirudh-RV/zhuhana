# To Run Cortex

> docker compose -f docker-compose/docker-compose-local.yaml up -d --build
> docker compose -f docker-compose/docker-compose-local.yaml up --build

# To stop

> docker compose -f docker-compose/docker-compose-local.yaml down

# To run Ollama

> docker compose -f docker-compose/docker-compose-ollama-local.yaml up -d --build

# To Migrate

DATABASE_URL=postgres://cortex:password@localhost:5433/cortex sqlx migrate run

export DATABASE_URL=postgres://cortex:password@localhost:5425/cortex
cargo sqlx prepare --check

export DATABASE_URL=postgres://cortex:password@localhost:5425/cortex
cargo build
