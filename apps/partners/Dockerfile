# ===========================================
# Stage 1: Dependencies (cached separately)
# ===========================================
FROM node:23-alpine AS deps

WORKDIR /app

# Copy only package files for dependency caching
COPY package*.json ./

# Install dependencies (cached if package*.json unchanged)
RUN npm ci

# ===========================================
# Stage 2: Builder
# ===========================================
FROM node:23-alpine AS builder

WORKDIR /app

# Copy cached dependencies
COPY --from=deps /app/node_modules ./node_modules

# Copy source code
COPY . .

# Build arguments for environment variables
ARG VITE_API_URL
ARG VITE_GRAPHQL_URL

ENV VITE_API_URL=$VITE_API_URL
ENV VITE_GRAPHQL_URL=$VITE_GRAPHQL_URL

# Build the application
RUN npm run build

# ===========================================
# Stage 3: Production
# ===========================================
FROM nginx:alpine

# Copy nginx configuration template
COPY nginx.conf /etc/nginx/templates/default.conf.template

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Default port
ENV PORT=80

# Expose port (will be overridden by PORT env var)
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q --spider http://localhost:${PORT}/ || exit 1

# Nginx will auto-substitute env vars in templates
CMD ["nginx", "-g", "daemon off;"]
