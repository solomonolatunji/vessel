#!/usr/bin/env bash
set -eo pipefail

CODEDOCK_DIR=${CODEDOCK_DIR:-/codedock}

if [ ! -d "$CODEDOCK_DIR/data" ]; then
  if [ -d "./data" ]; then
    CODEDOCK_DIR="."
  else
    echo "❌ No Codedock data directory found at $CODEDOCK_DIR/data."
    exit 1
  fi
fi

BACKUP_DIR="${CODEDOCK_DIR}/data/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/codedock-backup-${TIMESTAMP}.tar.gz"

echo "📦 Starting Codedock automated backup to ${BACKUP_FILE}..."
mkdir -p "${BACKUP_DIR}"

# Archive data folder (SQLite DB, Traefik configs) excluding existing backups
tar --exclude="backups" -czf "${BACKUP_FILE}" -C "${CODEDOCK_DIR}" data
echo "✅ Backup created successfully: ${BACKUP_FILE}"
