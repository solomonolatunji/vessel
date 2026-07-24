# Stage 1: Build the Dashboard GUI (TanStack + Vite)
FROM node:22-alpine AS dashboard-builder
WORKDIR /app

# Copy root monorepo files and workspaces
COPY package*.json ./
COPY apps/dashboard/package*.json ./apps/dashboard/
COPY apps/web/package*.json ./apps/web/

# Install dependencies
RUN npm ci

# Copy dashboard source code and build
COPY apps/dashboard/ ./apps/dashboard/
RUN npm run build:dashboard

# Stage 2: Build the static Go daemon (`codedockd`)
FROM golang:1.25-alpine AS daemon-builder
WORKDIR /src

# Install git and certificates
RUN apk add --no-cache git ca-certificates tzdata

# Copy Go modules manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source code
COPY . .

# Copy built dashboard static assets to be embedded
COPY --from=dashboard-builder /app/apps/dashboard/dist ./apps/dashboard/dist

# Accept version via build arguments (defaults to dev)
ARG VERSION=dev

# Build self-contained binary with CGO disabled and inject the version
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s -X main.codedockVersion=${VERSION}" -o /codedockd ./cmd/codedockd

# Stage 3: Minimal Production Runtime
FROM alpine:3.21 AS production
WORKDIR /var/www/codedock

# Install ca-certificates, docker-cli, git, and openssh-client for container orchestration and git cloning
RUN apk add --no-cache ca-certificates tzdata docker-cli git openssh-client curl

# Copy binary from daemon-builder
COPY --from=daemon-builder /codedockd /usr/local/bin/codedockd

# Ensure data directory exists
RUN mkdir -p /var/www/codedock/data

# Environment variables
ENV PORT=8080 \
    CODEDOCK_DATA_DIR=/var/www/codedock/data

EXPOSE 8080 80 443

VOLUME ["/var/www/codedock/data"]

ENTRYPOINT ["codedockd"]
