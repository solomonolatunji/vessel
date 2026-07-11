#!/usr/bin/env bash
# Vessel 1-Click Installer
# Usage: curl -fsSL https://get.vessel.dev | sh
set -eo pipefail

RELEASE=${VESSEL_VERSION:-latest}
VESSEL_DIR=/vessel
COMPOSE_URL="https://raw.githubusercontent.com/vesslhq/vessl/main/docker-compose.yml"

echo "🛰️  Vessel — Installing v${RELEASE}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# --- Root check ---
if [ "$EUID" -ne 0 ]; then
  echo "❌ Please run as root (or with sudo)."
  exit 1
fi

# --- Docker ---
if ! command -v docker &> /dev/null; then
  echo "📦 Installing Docker..."
  curl -fsSL https://get.docker.com | sh
  systemctl enable --now docker
fi

if ! docker info &> /dev/null; then
  echo "⏳ Waiting for Docker..."
  sleep 3
fi

# --- Directory setup ---
mkdir -p "$VESSEL_DIR"/data/{backups,caddy,builds}

# --- Pull files ---
echo "⬇️  Fetching configuration files..."
curl -fsSL "$COMPOSE_URL" -o "$VESSEL_DIR/docker-compose.yml"

if [ ! -f "$VESSEL_DIR/.env" ]; then
  echo "🔑 Generating .env file..."
  ENV_URL="https://raw.githubusercontent.com/vesslhq/vessl/main/.env.example"
  curl -fsSL "$ENV_URL" -o "$VESSEL_DIR/.env"
  # Generate a random 32-character string for JWT secret
  JWT_SECRET=$(head -c 24 /dev/urandom | base64)
  sed -i "s/VESSEL_JWT_SECRET=.*/VESSEL_JWT_SECRET=${JWT_SECRET}/" "$VESSEL_DIR/.env"
fi

# --- Pull & start ---
echo "🐳 Pulling vessel:v${RELEASE}..."
docker compose -f "$VESSEL_DIR/docker-compose.yml" pull
echo "🚀 Starting Vessel..."
docker compose -f "$VESSEL_DIR/docker-compose.yml" up -d

# --- systemd service (optional) ---
if command -v systemctl &> /dev/null; then
  cat > /etc/systemd/system/vessel.service <<'SERVICE'
[Unit]
Description=Vessel – Self-hosted PaaS
After=docker.service
Requires=docker.service

[Service]
Restart=always
RestartSec=10
WorkingDirectory=/vessel
ExecStartPre=-/usr/bin/docker compose -f /vessel/docker-compose.yml down
ExecStart=/usr/bin/docker compose -f /vessel/docker-compose.yml up
ExecStop=/usr/bin/docker compose -f /vessel/docker-compose.yml down

[Install]
WantedBy=multi-user.target
SERVICE
  systemctl daemon-reload
  systemctl enable --now vessel.service
fi

echo "✅ Vessel is running! Access the dashboard at http://$(curl -4fsS ifconfig.me 2>/dev/null || echo 'your-server-ip'):3000"
echo "📖 Docs: https://docs.vessel.dev"
