.PHONY: help build build-backend build-frontend run run-backend run-frontend test test-backend test-frontend clean docker-build docker-run docker-stop docker-down docker-logs docker-ps install install-backend install-frontend

# Default target
help:
	@echo "Available targets:"
	@echo "  make install          - Install dependencies for both frontend and backend"
	@echo "  make build            - Build both frontend and backend"
	@echo "  make run              - Run both frontend and backend"
	@echo "  make test             - Run tests for both frontend and backend"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-run       - Run Docker container"
	@echo "  make docker-stop      - Stop Docker container"
	@echo "  make docker-down      - Stop and remove Docker containers"
	@echo "  make docker-logs      - View Docker logs"
	@echo "  make docker-ps        - List running Docker containers"
	@echo "  make docker-compose-up - Start services with docker-compose"
	@echo "  make docker-compose-down - Stop services with docker-compose"

# Install dependencies
install: install-backend install-frontend

install-backend:
	@echo "Installing backend dependencies..."
	cd backend && go mod download

install-frontend:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Build targets
build: build-backend build-frontend

build-backend:
	@echo "Building backend..."
	cd backend && go build -o bin/server main.go

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm run build

# Run targets
run: run-backend run-frontend

run-backend:
	@echo "Running backend..."
	cd backend && go run main.go

run-frontend:
	@echo "Running frontend..."
	cd frontend && npm run dev

# Test targets
test: test-backend test-frontend

test-backend:
	@echo "Running backend tests..."
	cd backend && go test ./...

test-frontend:
	@echo "Running frontend tests..."
	cd frontend && npm run lint

# Clean targets
clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules/.vite

# Docker targets
docker-build:
	@echo "Building Docker image..."
	docker build -t tcmp-fin-app:latest .

docker-run:
	@echo "Running Docker container..."
	docker run -d -p 8080:8080 --name tcmp-fin-app tcmp-fin-app:latest

docker-stop:
	@echo "Stopping Docker container..."
	docker stop tcmp-fin-app || true

docker-down:
	@echo "Stopping and removing Docker container..."
	docker stop tcmp-fin-app || true
	docker rm tcmp-fin-app || true

docker-logs:
	@echo "Viewing Docker logs..."
	docker logs -f tcmp-fin-app

docker-ps:
	@echo "Listing Docker containers..."
	docker ps -a | grep tcmp-fin-app || echo "No containers found"

# Docker Compose targets
docker-compose-up:
	@echo "Starting services with docker-compose..."
	docker-compose up -d

docker-compose-down:
	@echo "Stopping services with docker-compose..."
	docker-compose down

docker-compose-logs:
	@echo "Viewing docker-compose logs..."
	docker-compose logs -f

docker-compose-build:
	@echo "Building images with docker-compose..."
	docker-compose build

