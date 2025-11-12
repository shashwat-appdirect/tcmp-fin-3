# Stage 1: Build Frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package.json ./

# Install dependencies
RUN npm install

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend

# Install git (needed for some Go dependencies)
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum* ./

# Download dependencies
RUN go mod download

# Copy backend source
COPY backend/ ./

# Build backend binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 3: Final Image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy backend binary from builder
COPY --from=backend-builder /app/backend/main .

# Copy frontend build from frontend-builder
COPY --from=frontend-builder /app/frontend/dist ./static

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]

