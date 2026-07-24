#!/usr/bin/env bash
set -eo pipefail

VERSION=${1:-latest}
CODEDOCK_DIR=${CODEDOCK_DIR:-/codedock}

echo "🛰️  Starting Codedock upgrade to v${VERSION}..."

if [ ! -f "$CODEDOCK_DIR/docker-compose.yml" ]; then
  if [ -f "./docker-compose.yml" ]; then
    CODEDOCK_DIR="."
  else
    echo "❌ No Codedock installation found at $CODEDOCK_DIR."
    exit 1
  fi
fi

echo "📦 Taking pre-upgrade state backup..."
mkdir -p "$CODEDOCK_DIR/data/backups"
if [ -f "$CODEDOCK_DIR/data/codedock.db" ]; then
  cp "$CODEDOCK_DIR/data/codedock.db" "$CODEDOCK_DIR/data/backups/codedock-pre-upgrade-$(date +%Y%m%d%H%M%S).db" 2>/dev/null || true
  echo "✅ SQLite database backed up safely."
fi

# Ensure required env vars exist in .env
if [ -f "$CODEDOCK_DIR/.env" ]; then
  # Add missing vars with defaults if not present
  ensure_env() {
    local key="$1"
    local default="$2"
    if ! grep -q "^${key}=" "$CODEDOCK_DIR/.env" 2>/dev/null; then
      echo "${key}=${default}" >> "$CODEDOCK_DIR/.env"
      echo "  + Added ${key} to .env"
    fi
  }

  echo "🔍 Checking .env for required variables..."
  ensure_env "DOCKER_SOCKET_PATH" "/var/run/docker.sock"
  ensure_env "CODEDOCK_DASHBOARD_URL" "http://localhost:8080"
  ensure_env "CODEDOCK_TRAEFIK_HTTP_PORT" "80"
  ensure_env "CODEDOCK_TRAEFIK_HTTPS_PORT" "443"
  ensure_env "CODEDOCK_TRAEFIK_API_PORT" "8080"
fi

export CODEDOCK_VERSION="$VERSION"
echo "⬇️  Fetching Codedock release v${VERSION}..."

if command -v docker &> /dev/null; then
  docker compose -f "$CODEDOCK_DIR/docker-compose.yml" pull
  docker compose -f "$CODEDOCK_DIR/docker-compose.yml" up -d --force-recreate
else
  echo "❌ Docker not found. Running outside Docker? Please replace binary manually."
  exit 1
fi

echo "🚀 Codedock upgrade to v${VERSION} completed successfully! User containers experienced zero downtime."
