# Stage 1: Build the Dashboard GUI (TanStack + Vite)
FROM node:22-alpine AS dashboard-builder
WORKDIR /app

# Copy root monorepo files and workspaces
COPY package*.json ./
COPY dashboard/package*.json ./dashboard/
COPY website/package*.json ./website/

# Install dependencies
RUN npm ci

# Copy dashboard source code and build
COPY dashboard/ ./dashboard/
RUN npm run build:dashboard

# Stage 2: Build the static Go daemon (`vesseld`)
FROM golang:1.24-alpine AS daemon-builder
WORKDIR /src

# Install git and certificates
RUN apk add --no-cache git ca-certificates tzdata

# Copy Go modules manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source code
COPY . .

# Copy built dashboard static assets to be embedded
COPY --from=dashboard-builder /app/dashboard/dist ./dashboard/dist

# Build self-contained binary with CGO disabled
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /vesseld ./cmd/vesseld

# Stage 3: Minimal Production Runtime
FROM alpine:3.21 AS production
WORKDIR /var/www/vessel

# Install ca-certificates, docker-cli, git, and openssh-client for container orchestration and git cloning
RUN apk add --no-cache ca-certificates tzdata docker-cli git openssh-client curl

# Copy binary from daemon-builder
COPY --from=daemon-builder /vesseld /usr/local/bin/vesseld

# Ensure data directory exists
RUN mkdir -p /var/www/vessel/data

# Environment variables
ENV PORT=8080 \
    VESSEL_DATA_DIR=/var/www/vessel/data

EXPOSE 8080 80 443

VOLUME ["/var/www/vessel/data"]

ENTRYPOINT ["vesseld"]
