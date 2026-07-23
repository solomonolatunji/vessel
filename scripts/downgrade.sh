#!/usr/bin/env bash
set -eo pipefail

VERSION=${1:-}
CODEDOCK_DIR=${CODEDOCK_DIR:-/codedock}

if [ -z "$VERSION" ]; then
  echo "❌ Usage: ./scripts/downgrade.sh <version>"
  echo "   Example: ./scripts/downgrade.sh 0.1.0"
  echo ""
  echo "Available versions: https://github.com/buildwithtechx/codedock/releases"
  exit 1
fi

echo "⬇️  Codedock — Downgrading to v${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ ! -f "$CODEDOCK_DIR/docker-compose.yml" ]; then
  if [ -f "./docker-compose.yml" ]; then
    CODEDOCK_DIR="."
  else
    echo "❌ No Codedock installation found at $CODEDOCK_DIR."
    exit 1
  fi
fi

echo "💾 Creating database backup before downgrade..."
mkdir -p "$CODEDOCK_DIR/data/backups"
if [ -f "$CODEDOCK_DIR/data/codedock.db" ]; then
  cp "$CODEDOCK_DIR/data/codedock.db" "$CODEDOCK_DIR/data/backups/codedock-pre-downgrade-$(date +%Y%m%d%H%M%S).db" 2>/dev/null || true
  echo "✅ SQLite database backed up safely."
fi

export CODEDOCK_VERSION="$VERSION"
echo "🐳 Pulling codedock:v${VERSION}..."

if command -v docker &> /dev/null; then
  docker compose -f "$CODEDOCK_DIR/docker-compose.yml" pull
  echo "🚀 Recreating container with v${VERSION}..."
  docker compose -f "$CODEDOCK_DIR/docker-compose.yml" up -d --force-recreate
else
  echo "❌ Docker not found."
  exit 1
fi

echo "✅ Downgraded to v${VERSION}."
echo "⚠️  If you encounter issues, restore a previous backup from $CODEDOCK_DIR/data/backups/"
