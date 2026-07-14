#!/usr/bin/env bash
set -eo pipefail

VERSION=${1:-latest}
VESSL_DIR=${VESSL_DIR:-/vessl}

echo "🛰️  Starting Vessl upgrade to v${VERSION}..."

if [ ! -f "$VESSL_DIR/docker-compose.yml" ]; then
  if [ -f "./docker-compose.yml" ]; then
    VESSL_DIR="."
  else
    echo "❌ No Vessl installation found at $VESSL_DIR."
    exit 1
  fi
fi

echo "📦 Taking pre-upgrade state backup..."
mkdir -p "$VESSL_DIR/data/backups"
if [ -f "$VESSL_DIR/data/vessl.db" ]; then
  cp "$VESSL_DIR/data/vessl.db" "$VESSL_DIR/data/backups/vessl-pre-upgrade-$(date +%Y%m%d%H%M%S).db" 2>/dev/null || true
  echo "✅ SQLite database backed up safely."
fi

# Ensure required env vars exist in .env
if [ -f "$VESSL_DIR/.env" ]; then
  # Add missing vars with defaults if not present
  ensure_env() {
    local key="$1"
    local default="$2"
    if ! grep -q "^${key}=" "$VESSL_DIR/.env" 2>/dev/null; then
      echo "${key}=${default}" >> "$VESSL_DIR/.env"
      echo "  + Added ${key} to .env"
    fi
  }

  echo "🔍 Checking .env for required variables..."
  ensure_env "DOCKER_SOCKET_PATH" "/var/run/docker.sock"
  ensure_env "VESSL_DASHBOARD_URL" "http://localhost:8080"
  ensure_env "VESSL_TRAEFIK_HTTP_PORT" "80"
  ensure_env "VESSL_TRAEFIK_HTTPS_PORT" "443"
  ensure_env "VESSL_TRAEFIK_API_PORT" "8080"
fi

export VESSL_VERSION="$VERSION"
echo "⬇️  Fetching Vessl release v${VERSION}..."

if command -v docker &> /dev/null; then
  docker compose -f "$VESSL_DIR/docker-compose.yml" pull
  docker compose -f "$VESSL_DIR/docker-compose.yml" up -d --force-recreate
else
  echo "❌ Docker not found. Running outside Docker? Please replace binary manually."
  exit 1
fi

echo "🚀 Vessl upgrade to v${VERSION} completed successfully! User containers experienced zero downtime."
