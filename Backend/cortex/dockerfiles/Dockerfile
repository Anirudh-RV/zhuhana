FROM rust:1.82-bookworm as builder

WORKDIR /app
COPY . .

RUN cargo fetch
RUN cargo build --release

FROM debian:bookworm-slim

# Install OpenSSL 3
RUN apt-get update && apt-get install -y \
    libssl3 \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/target/release/cortex /usr/local/bin/llmservice
EXPOSE 3000
CMD ["llmservice"]
