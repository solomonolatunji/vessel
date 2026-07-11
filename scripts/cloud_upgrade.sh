#!/usr/bin/env bash
set -eo pipefail

echo "☁️  Vessel Cloud — Zero-Downtime Upgrade"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

VESSEL_DIR=${VESSEL_DIR:-/opt/vessel-cloud}

if [ -d "$VESSEL_DIR" ]; then
  cd "$VESSEL_DIR"
else
  echo "❌ Cloud directory not found: $VESSEL_DIR"
  exit 1
fi

echo "💾 Backing up Cloud Postgres DB..."
mkdir -p ./backups
docker exec vessel-postgres pg_dump -U vessel -d vessel_cloud > ./backups/cloud-$(date +%s).sql

echo "🐳 Pulling latest cloud images..."
docker compose pull

echo "🚀 Restarting cloud API with zero-downtime rolling update..."
docker compose up -d --no-deps --build vessel-cloud-api
docker exec vessel-cloud-api /vesseld --migrate

echo "✅ Cloud Upgrade Complete."
