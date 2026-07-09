#!/usr/bin/env bash
# Vessel Downgrade
# Usage: curl -fsSL https://get.vessel.dev/downgrade | sh -s <version>
set -eo pipefail

VERSION=${1:-}
VESSEL_DIR=${VESSEL_DIR:-/vessel}

if [ -z "$VERSION" ]; then
  echo "❌ Usage: $0 <version>"
  echo "   Example: $0 0.1.0"
  echo ""
  echo "Available versions: https://github.com/solomonolatunji/vessel/releases"
  exit 1
fi

echo "⬇️  Vessel — Downgrading to v${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ ! -f "$VESSEL_DIR/docker-compose.yml" ]; then
  echo "❌ No Vessel installation found at $VESSEL_DIR."
  exit 1
fi

echo "💾 Creating database backup before downgrade..."
mkdir -p "$VESSEL_DIR/data/backups"
docker compose -f "$VESSEL_DIR/docker-compose.yml" exec -T vesseld cat /data/vessel.db > "$VESSEL_DIR/data/backups/vessel-$(date +%Y%m%d%H%M%S).db" 2>/dev/null || \
  cp "$VESSEL_DIR/data/vessel.db" "$VESSEL_DIR/data/backups/vessel-$(date +%Y%m%d%H%M%S).db" 2>/dev/null || \
  echo "⚠️  Could not backup database. Proceeding anyway."

export VESSEL_VERSION="$VERSION"
echo "🐳 Pulling vessel:v${VERSION}..."
docker compose -f "$VESSEL_DIR/docker-compose.yml" pull
echo "🚀 Recreating container with v${VERSION}..."
docker compose -f "$VESSEL_DIR/docker-compose.yml" up -d --force-recreate

echo "✅ Downgraded to v${VERSION}."
echo "⚠️  If you encounter issues, restore a previous backup from $VESSEL_DIR/data/backups/"
