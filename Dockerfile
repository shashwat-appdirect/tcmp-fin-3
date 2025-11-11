# Multi-stage Dockerfile for building both frontend and backend

# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./

# Install dependencies
RUN npm ci

# Copy frontend source code
COPY frontend/ ./

# Build frontend
RUN npm run build

# Stage 2: Build Backend
# Note: If go.mod specifies go 1.25.0, you may need to adjust this tag
# Use golang:alpine for latest, or golang:1.23-alpine for a specific version
FROM golang:alpine AS backend-builder

WORKDIR /app/backend

# Install git (needed for some Go dependencies)
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Build backend binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server main.go

# Stage 3: Final Runtime Image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates wget

# Copy backend binary from builder
COPY --from=backend-builder /app/backend/server ./server

# Copy frontend build from builder
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Note: service-account.json is NOT copied into the image for security
# It should be mounted as a volume at runtime (see docker-compose.yml)

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ARG ADMIN_PASSWORD=${ADMIN_PASSWORD:-admin123}
ENV ADMIN_PASSWORD=${ADMIN_PASSWORD}
ENV SERVICE_ACCOUNT_PATH=/app/service-account.json
ENV FRONTEND_DIR=/app/frontend/dist
# Note: For Google Cloud Run deployment:
# - GOOGLE_CLOUD_PROJECT is automatically set by Cloud Run
# - CLIENT_ID must be set as an environment variable in Cloud Run (this is your custom client identifier)
# - The service account file is not needed when using Application Default Credentials (ADC)

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/attendees/count || exit 1

# Run backend server
CMD ["./server"]

