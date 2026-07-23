SERVICES = dashboard web docs
BINARY_NAME = codedockd
BUILD_DIR = bin

.PHONY: help all build build-daemon build-dashboard dev dev-dryrun dev-daemon dev-dashboard clean check fmt test docker-build docker-up docker-down

all: check build

help:
	@echo "Codedock — available commands:"
	@echo ""
	@echo "  make check             Run Go fmt + vet"
	@echo "  make fmt               Run Go fmt only"
	@echo "  make test              Run Go tests"
	@echo "  make build             Build daemon binary + dashboard"
	@echo "  make dev               Run daemon + dashboard concurrently"
	@echo "  make dev-dryrun        Run daemon (with DEPLOY_DRY_RUN=true) + dashboard concurrently"
	@echo "  make dev-daemon        Run Go daemon in dev mode"
	@echo "  make dev-dashboard     Run dashboard dev server"
	@echo "  make clean             Remove build artifacts"
	@echo "  make docker-build      Build Docker image"
	@echo "  make docker-up         Start Docker stack"
	@echo "  make docker-down       Stop Docker stack"


check:
	@echo "🔍 Running Go checks and formatting..."
	go fmt ./...
	go vet ./...  

fmt:
	@echo "🔍 Formatting Go code..."
	go fmt ./...

test:
	@echo "🧪 Running full test suite..."
	go test ./... -v

build: build-dashboard build-daemon
	@echo "✅ Build complete! Binaries available in $(BUILD_DIR)/ and GUI at dashboard/dist"

build-daemon:
	@echo "⚙️  Building Go daemon binary ($(BINARY_NAME))..."
	mkdir -p $(BUILD_DIR)
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/codedockd

build-dashboard:
	@echo "💻 Building TanStack + Vite Dashboard GUI..."
	npm run build:dashboard

dev:
	@echo "🚀 Launching self-hosted backend daemon and frontend GUI concurrently..."
	npx concurrently -k "make dev-daemon" "make dev-dashboard"

dev-dryrun:
	@echo "🚀 Launching with DEPLOY_DRY_RUN=true and frontend GUI concurrently..."
	npx concurrently -k "DEPLOY_DRY_RUN=true make dev-daemon" "make dev-dashboard"

dev-daemon:
	@echo "🚀 Running Go daemon in dev mode with live reload..."
	go run github.com/air-verse/air@latest -c .air.toml

dev-dashboard:
	@echo "💻 Running Dashboard dev server on port 3000..."
	npm run dev:dashboard

dev-website:
	@echo "🌐 Running Astro Marketing site dev server..."
	npm run dev:website

docker-build:
	@echo "🐳 Building Docker image..."
	docker compose build

docker-up:
	@echo "🐳 Starting Codedock via Docker Compose..."
	docker compose up -d

docker-down:
	@echo "🐳 Stopping Codedock Docker stack..."
	docker compose down

clean:
	@echo "🧹 Cleaning builds and temporary binaries..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)

