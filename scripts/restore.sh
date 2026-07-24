#!/usr/bin/env bash
set -eo pipefail

if [ -z "$1" ]; then
  echo "❌ Usage: ./scripts/restore.sh <path-to-backup.tar.gz>"
  exit 1
fi

BACKUP_FILE="$1"
CODEDOCK_DIR=${CODEDOCK_DIR:-/codedock}

if [ ! -f "${BACKUP_FILE}" ]; then
  echo "❌ Error: Backup file ${BACKUP_FILE} not found!"
  exit 1
fi

if [ ! -d "$CODEDOCK_DIR/data" ]; then
  if [ -d "./data" ]; then
    CODEDOCK_DIR="."
  else
    echo "❌ No Codedock data directory found at $CODEDOCK_DIR/data."
    exit 1
  fi
fi

echo "⚠️  Restoring Codedock state from ${BACKUP_FILE} into ${CODEDOCK_DIR}/data..."
# Extract relative to CODEDOCK_DIR
tar -xzf "${BACKUP_FILE}" -C "${CODEDOCK_DIR}"

# Restarting the service
echo "🔄 Restarting Codedock container to apply restored data..."
if command -v docker &> /dev/null && [ -f "$CODEDOCK_DIR/docker-compose.yml" ]; then
  docker compose -f "$CODEDOCK_DIR/docker-compose.yml" restart codedock
fi

echo "✅ Restore completed successfully! Codedock is back online."
