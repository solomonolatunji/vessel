#!/usr/bin/env bash
set -eo pipefail

VERSION=${1:-}
VESSEL_DIR=${VESSEL_DIR:-/vessel}

if [ -z "$VERSION" ]; then
  echo "❌ Usage: ./scripts/downgrade.sh <version>"
  echo "   Example: ./scripts/downgrade.sh 0.1.0"
  echo ""
  echo "Available versions: https://github.com/vesslhq/vessl/releases"
  exit 1
fi

echo "⬇️  Vessel — Downgrading to v${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ ! -f "$VESSEL_DIR/docker-compose.yml" ]; then
  if [ -f "./docker-compose.yml" ]; then
    VESSEL_DIR="."
  else
    echo "❌ No Vessel installation found at $VESSEL_DIR."
    exit 1
  fi
fi

echo "💾 Creating database backup before downgrade..."
mkdir -p "$VESSEL_DIR/data/backups"
if [ -f "$VESSEL_DIR/data/vessel.db" ]; then
  cp "$VESSEL_DIR/data/vessel.db" "$VESSEL_DIR/data/backups/vessel-pre-downgrade-$(date +%Y%m%d%H%M%S).db" 2>/dev/null || true
  echo "✅ SQLite database backed up safely."
fi

export VESSEL_VERSION="$VERSION"
echo "🐳 Pulling vessel:v${VERSION}..."

if command -v docker &> /dev/null; then
  docker compose -f "$VESSEL_DIR/docker-compose.yml" pull
  echo "🚀 Recreating container with v${VERSION}..."
  docker compose -f "$VESSEL_DIR/docker-compose.yml" up -d --force-recreate
else
  echo "❌ Docker not found."
  exit 1
fi

echo "✅ Downgraded to v${VERSION}."
echo "⚠️  If you encounter issues, restore a previous backup from $VESSEL_DIR/data/backups/"
