.PHONY: build run test docker-build docker-run clean

# Environment Variables
SERVICE_NAME=service-availability
DOCKER_IMAGE=$(SERVICE_NAME):latest
DOCKER_COMPOSE_FILE=infrastructure/docker-compose.yml

# Build the Go application
build:
	@echo "🔨 Building $(SERVICE_NAME)..."
	go build -o bin/$(SERVICE_NAME) main.go

# Run the application locally
run: build
	@echo "🚀 Running $(SERVICE_NAME)..."
	./bin/$(SERVICE_NAME)

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./tests -v

# Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_IMAGE) -f infrastructure/Dockerfile .

# Run Docker Compose
docker-run:
	@echo "📦 Running with Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

# Clean up
clean:
	@echo "🧹 Cleaning up..."
	rm -rf bin/ $(SERVICE_NAME) coverage.out
