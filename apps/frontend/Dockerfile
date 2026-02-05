# Stage 1: Build the application
FROM node:22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy frontend package files
COPY package*.json ./frontend/

# Install dependencies
WORKDIR /app/frontend
RUN npm ci

# Set build arguments
ARG VITE_API_BASE_URL
ARG VITE_ZITADEL_AUTHORITY
ARG VITE_ZITADEL_CLIENT_ID

ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
ENV VITE_ZITADEL_AUTHORITY=$VITE_ZITADEL_AUTHORITY
ENV VITE_ZITADEL_CLIENT_ID=$VITE_ZITADEL_CLIENT_ID

# Copy frontend source code
COPY . .

# Build the application
RUN npm run build

# Stage 2: Production image with Nginx
FROM nginx:alpine

# Copy built assets from builder stage
COPY --from=builder /app/frontend/dist /usr/share/nginx/html

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://127.0.0.1/health || exit 1

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
