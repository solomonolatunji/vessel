SERVICES = dashboard cloud web docs
BINARY_NAME = vesseld
BUILD_DIR = bin

.PHONY: all build build-daemon build-dashboard dev dev-daemon dev-dashboard clean check fmt test install update organize-imports deploy deploy-web deploy-cloud deploy-docs docker-build docker-up docker-down

all: check build

install:
	@echo "📦 Installing all dependencies..."
	npm install

update:
	@echo "⬆️  Checking for dependency updates..."
	@for dir in . $(SERVICES); do \
		if [ -f "$$dir/package.json" ]; then \
			echo "  📦 $$dir"; \
			(cd $$dir && npx npm-check-updates -u 2>/dev/null) || true; \
		fi \
	done
	@echo "✅ Update check complete. Run 'make install' to apply."

organize-imports:
	@echo "🗂️  Organizing TypeScript imports..."
	@for dir in $(SERVICES); do \
		if [ -f "$$dir/package.json" ]; then \
			echo "  📦 $$dir"; \
			(cd $$dir && npx organize-imports-cli tsconfig*.json 2>/dev/null) || true; \
		fi \
	done

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
	@echo "✅ Build complete! Binary available at $(BUILD_DIR)/$(BINARY_NAME) and GUI at dashboard/dist"

build-daemon:
	@echo "⚙️  Building Go daemon binary ($(BINARY_NAME))..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/vesseld

build-dashboard:
	@echo "💻 Building TanStack + Vite Dashboard GUI..."
	npm run build:dashboard

dev:
	@echo "🚀 Launching backend daemon and frontend GUI concurrently..."
	npx concurrently -k "make dev-daemon" "make dev-dashboard"

dev-daemon:
	@echo "🚀 Running Go daemon in dev mode..."
	go run ./cmd/vesseld

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
	@echo "🐳 Starting Vessel via Docker Compose..."
	docker compose up -d

docker-down:
	@echo "🐳 Stopping Vessel Docker stack..."
	docker compose down

clean:
	@echo "🧹 Cleaning builds and temporary binaries..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)

deploy-web:
	@echo "🌐 Deploying web to Cloudflare Pages..."
	npm run deploy:web

deploy-cloud:
	@echo "☁️  Deploying cloud to Cloudflare Pages..."
	npm run deploy:cloud

deploy-docs:
	@echo "📖 Deploying docs to Cloudflare Pages..."
	npm run deploy:docs

deploy:
	@echo "🚀 Deploying all frontends to Cloudflare Pages..."
	npm run deploy:all
