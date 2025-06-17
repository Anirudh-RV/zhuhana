# Stage 1: Build React app
FROM node:20-alpine as builder

WORKDIR /app

# Copy and build frontend
COPY frontend ./frontend
WORKDIR /app/frontend
RUN npm install && npm run build

# Stage 2: Build Node.js server
FROM node:20-alpine

WORKDIR /app

# Copy backend code
COPY package*.json tsconfig.json ./
COPY src ./src

# Install backend deps
RUN npm install

# Copy built frontend from builder
COPY --from=builder /app/frontend/build ./frontend/build

# Build TypeScript
RUN npx tsc

# Expose port
EXPOSE 3000

CMD ["node", "dist/index.js"]
