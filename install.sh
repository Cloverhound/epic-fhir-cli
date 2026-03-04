#!/bin/sh
set -e

# FHIR CLI installer
# Usage: curl -fsSL https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/install.sh | sh

REPO="Cloverhound/epic-fhir-cli"
BINARY="fhir-cli"
INSTALL_DIR="${HOME}/.local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  darwin) OS="darwin" ;;
  linux)  OS="linux" ;;
  *)      echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64)  ARCH="arm64" ;;
  *)              echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version
echo "Fetching latest release..."
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
if [ -z "$VERSION" ]; then
  echo "Error: could not determine latest version"
  exit 1
fi
echo "Latest version: v${VERSION}"

# Download
TARBALL="${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${TARBALL}"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${URL}..."
curl -fsSL "$URL" -o "${TMPDIR}/${TARBALL}"

# Extract
tar -xzf "${TMPDIR}/${TARBALL}" -C "$TMPDIR"

# Install
mkdir -p "$INSTALL_DIR"
mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
chmod +x "${INSTALL_DIR}/${BINARY}"

echo "Installed ${BINARY} v${VERSION} to ${INSTALL_DIR}/${BINARY}"
echo ""

# Warn if INSTALL_DIR is not on PATH
case ":${PATH}:" in
  *":${INSTALL_DIR}:"*) ;;
  *)
    # Detect shell config file
    SHELL_NAME=$(basename "${SHELL:-/bin/sh}")
    case "$SHELL_NAME" in
      zsh)  RC_FILE=~/.zshrc ;;
      bash) RC_FILE=~/.bashrc ;;
      fish) RC_FILE=~/.config/fish/config.fish ;;
      *)    RC_FILE=~/.profile ;;
    esac

    echo "Note: ${INSTALL_DIR} is not in your PATH. Run:"
    echo ""
    if [ "$SHELL_NAME" = "fish" ]; then
      echo "  fish_add_path ${INSTALL_DIR}"
    else
      echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ${RC_FILE} && source ${RC_FILE}"
    fi
    echo ""
    ;;
esac

echo "Get started:"
echo "  fhir-cli config init    # Interactive setup"
echo "  fhir-cli auth token     # Get an access token"
echo "  fhir-cli patient search --family Smith"
