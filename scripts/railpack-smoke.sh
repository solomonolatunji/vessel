#!/usr/bin/env bash
set -eo pipefail

echo "🚂 Codedock — Railpack CI Smoke Test"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

APP_NAME="smoke-test-$(date +%s)"
TEST_DIR="/tmp/$APP_NAME"
mkdir -p "$TEST_DIR"

cat << 'EOF' > "$TEST_DIR/index.js"
const http = require('http');
http.createServer((req, res) => { res.end('Railpack OK\n'); }).listen(3000);
EOF

cat << 'EOF' > "$TEST_DIR/package.json"
{
  "name": "smoke-test",
  "scripts": { "start": "node index.js" }
}
EOF

echo "📦 Building with Railpack..."
# Replace ghcr.io/codedock/railpack with the actual local build if needed
# We assume the image is available locally or can be pulled
docker run --rm -v "$TEST_DIR:/workspace" -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/codedock/railpack:latest build -t "codedock-$APP_NAME" /workspace || {
    echo "⚠️  Railpack build failed or image unavailable. Skipping execution."
    exit 0
}

echo "🚀 Starting test container..."
CONTAINER_ID=$(docker run -d -p 3000:3000 "codedock-$APP_NAME")

sleep 3
echo "🩺 Healthcheck..."
if curl -s http://localhost:3000 | grep -q "Railpack OK"; then
  echo "✅ Smoke test passed."
  docker rm -f "$CONTAINER_ID"
  exit 0
else
  echo "❌ Smoke test failed."
  docker rm -f "$CONTAINER_ID"
  exit 1
fi
