#!/usr/bin/env bash
# codedock CLI Installer
# Usage: curl -fsSL https://get.codedock.dev/cli | sh
set -eo pipefail

REPO="buildwithtechx/codedock"
BINARY="codedock"
INSTALL_DIR="/usr/local/bin"
BOLD="\033[1m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
DIM="\033[2m"
NC="\033[0m"

echo -e "${BOLD}🛰️  codedock CLI Installer${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

detect_platform() {
  OS="$(uname -s)"
  ARCH="$(uname -m)"

  case "$OS" in
    Linux)  OS="linux" ;;
    Darwin) OS="darwin" ;;
    *)
      echo -e "${RED}❌ Unsupported OS: $OS${NC}"
      exit 1
      ;;
  esac

  case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    *)
      echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
      exit 1
      ;;
  esac

  echo "${OS}_${ARCH}"
}

get_latest_version() {
  curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null \
    | grep '"tag_name"' \
    | head -1 \
    | cut -d'"' -f4
}

install_via_go() {
  if command -v go &>/dev/null; then
    echo -e "${YELLOW}⚙️  No pre-built binary found. Installing via 'go install'...${NC}"
    go install "codedock.dev/codedock/cmd/codedock@latest"
    GOBIN=$(go env GOPATH)/bin
    if [ -f "$GOBIN/codedock" ]; then
      if [ -w "$INSTALL_DIR" ] || [ "$(id -u)" -eq 0 ]; then
        cp "$GOBIN/codedock" "$INSTALL_DIR/codedock"
        echo -e "${GREEN}✅ Installed via go install → $INSTALL_DIR/codedock${NC}"
      else
        echo -e "${YELLOW}⚠️  No write access to $INSTALL_DIR. Binary is at $GOBIN/codedock${NC}"
        echo -e "   Add \$(go env GOPATH)/bin to your PATH:"
        echo -e "   ${DIM}export PATH=\$PATH:\$(go env GOPATH)/bin${NC}"
      fi
    fi
    return 0
  fi
  return 1
}

PLATFORM=$(detect_platform)
echo -e "  Platform:  ${PLATFORM}"

TARGET_VERSION="${CODEDOCK_VERSION:-}"
if [ -z "$TARGET_VERSION" ]; then
  echo -e "  Version:   ${DIM}checking latest...${NC}"
  TARGET_VERSION=$(get_latest_version)
fi

if [ -z "$TARGET_VERSION" ]; then
  echo -e "${YELLOW}⚠️  Could not fetch latest release from GitHub.${NC}"
  install_via_go || {
    echo -e "${RED}❌ Could not install codedock. Install Go and run:${NC}"
    echo -e "   go install codedock.dev/codedock/cmd/codedock@latest"
    exit 1
  }
  exit 0
fi

echo -e "  Version:   ${TARGET_VERSION}"
echo ""

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TARGET_VERSION}/codedock_${PLATFORM}.tar.gz"

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo -e "${BOLD}⬇️  Downloading codedock ${TARGET_VERSION} (${PLATFORM})...${NC}"

if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/codedock.tar.gz" 2>/dev/null; then
  echo -e "${YELLOW}⚠️  Pre-built binary not available for this platform/version.${NC}"
  install_via_go || {
    echo -e "${RED}❌ Could not install codedock. Install Go and run:${NC}"
    echo -e "   go install codedock.dev/codedock/cmd/codedock@latest"
    exit 1
  }
  exit 0
fi

tar -xzf "$TMP_DIR/codedock.tar.gz" -C "$TMP_DIR"

BINARY_PATH="$TMP_DIR/codedock"
if [ ! -f "$BINARY_PATH" ]; then
  BINARY_PATH=$(find "$TMP_DIR" -name "codedock" -type f | head -1)
fi

if [ -z "$BINARY_PATH" ]; then
  echo -e "${RED}❌ Could not find codedock binary in archive.${NC}"
  exit 1
fi

chmod +x "$BINARY_PATH"

if [ -w "$INSTALL_DIR" ] || [ "$(id -u)" -eq 0 ]; then
  mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY"
  echo -e "${GREEN}✅ Installed → $INSTALL_DIR/$BINARY${NC}"
else
  LOCAL_BIN="$HOME/.local/bin"
  mkdir -p "$LOCAL_BIN"
  mv "$BINARY_PATH" "$LOCAL_BIN/$BINARY"
  echo -e "${GREEN}✅ Installed → $LOCAL_BIN/$BINARY${NC}"
  echo -e "${YELLOW}⚠️  $LOCAL_BIN is not in your PATH.${NC}"
  echo -e "   Add it:"
  echo -e "   ${DIM}export PATH=\$PATH:\$HOME/.local/bin${NC}"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ codedock ${TARGET_VERSION} installed successfully!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo -e "  ${BOLD}Get started:${NC}"
echo -e "  1. Connect to your self-hosted server:"
echo -e "     ${BOLD}codedock login${NC}"
echo -e ""
echo -e "  2. List your projects:"
echo -e "     ${BOLD}codedock project list${NC}"
echo -e ""
echo -e "  3. Trigger a deployment:"
echo -e "     ${BOLD}codedock deploy <service-id>${NC}"
echo ""
echo -e "  ${DIM}Docs: https://docs.codedock.dev/cli${NC}"
echo ""
