# Stage 1: Build the React app
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package.json and package-lock.json first to leverage Docker cache
COPY package*.json ./
RUN npm install

# Copy the rest of the app's source code
COPY . .

# Build the app for production
RUN npm run build

# Stage 2: Serve the app with Nginx
FROM nginx:stable-alpine

# Copy the built static files from the 'builder' stage
COPY --from=builder /app/build /usr/share/nginx/html

# Copy our custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80 (Nginx's default port)
EXPOSE 80

# Start Nginx when the container launches
CMD ["nginx", "-g", "daemon off;"]
