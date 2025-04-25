# Use a Node.js base image
FROM node:18-alpine

# Set the working directory
WORKDIR /

# Copy package.json and package-lock.json to install dependencies
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the TypeScript code
RUN npm run build

# Expose the port your application listens on
EXPOSE 3000

# Start the application
CMD ["node", "build/index.js"]
