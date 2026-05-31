#!/usr/bin/env bash
set -euo pipefail

APP_NAME="pdfsigner"
PACKAGE_NAME="pdfsigner"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILDROOT="${ROOT_DIR}/.debbuild"
DIST_DIR="${ROOT_DIR}/dist"

DESKTOP_FILE="${ROOT_DIR}/packaging/pdfsigner.desktop"
ICON_FILE="${ROOT_DIR}/packaging/pdfsigner.png"
BIN_FILE="${ROOT_DIR}/${APP_NAME}"

DEFAULT_VERSION="0.1.2"
DEFAULT_ARCH="amd64"

read -r -p "Package version [${DEFAULT_VERSION}]: " VERSION
VERSION="${VERSION:-$DEFAULT_VERSION}"

read -r -p "Debian architecture [${DEFAULT_ARCH}]: " DEB_ARCH
DEB_ARCH="${DEB_ARCH:-$DEFAULT_ARCH}"

echo "==> Checking project files"

for required_file in "${DESKTOP_FILE}" "${ICON_FILE}"; do
  if [[ ! -f "${required_file}" ]]; then
    echo "Error: ${required_file} not found"
    exit 1
  fi
done

echo "==> Cleaning previous build output"
rm -rf "${BUILDROOT}"
mkdir -p "${BUILDROOT}/DEBIAN" \
  "${BUILDROOT}/usr/bin" \
  "${BUILDROOT}/usr/share/applications" \
  "${BUILDROOT}/usr/share/icons/hicolor/256x256/apps" \
  "${BUILDROOT}/usr/share/doc/${PACKAGE_NAME}" \
  "${DIST_DIR}"

echo "==> Building Go binary"
cd "${ROOT_DIR}"
go mod download
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o "${BIN_FILE}" .

if [[ ! -f "${BIN_FILE}" ]]; then
  echo "Error: ${BIN_FILE} was not built"
  exit 1
fi

echo "==> Installing package files"
install -m 0755 "${BIN_FILE}" "${BUILDROOT}/usr/bin/pdfsigner"
install -m 0644 "${DESKTOP_FILE}" "${BUILDROOT}/usr/share/applications/pdfsigner.desktop"
install -m 0644 "${ICON_FILE}" "${BUILDROOT}/usr/share/icons/hicolor/256x256/apps/pdfsigner.png"

cat > "${BUILDROOT}/usr/share/doc/${PACKAGE_NAME}/README.Debian" <<'EOF'
pdfsigner requires CryptoPro CSP runtime tools:

- /opt/cprocsp/bin/amd64/certmgr
- /opt/cprocsp/bin/amd64/csptest

These tools are not declared as apt dependencies because CryptoPro CSP is
distributed separately.
EOF

cat > "${BUILDROOT}/DEBIAN/control" <<EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${DEB_ARCH}
Maintainer: shurshick <noreply@example.com>
Homepage: https://github.com/shurshick/pdf-signer
Description: PDF signer and visual electronic signature stamp tool
 Desktop application for adding a visual electronic signature stamp to PDF
 documents and working with certificate data from CryptoPro CSP.
EOF

echo "==> Building DEB"
dpkg-deb --build --root-owner-group "${BUILDROOT}" "${DIST_DIR}/${PACKAGE_NAME}_${VERSION}_${DEB_ARCH}.deb"

echo "==> Done"
echo "Version: ${VERSION}"
echo "Architecture: ${DEB_ARCH}"
echo "DEB packages:"
ls -lh "${DIST_DIR}"/*.deb
