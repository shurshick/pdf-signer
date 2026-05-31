#!/usr/bin/env bash
set -euo pipefail

APP_NAME="pdfsigner"
VERSION="0.1.0"
TOPDIR="${HOME}/rpmbuild"
WORKDIR="$(pwd)"
PKGROOT="${WORKDIR}/${APP_NAME}-${VERSION}"

rm -rf "${PKGROOT}"
mkdir -p "${PKGROOT}/packaging"

# Copy source tree, excluding common junk.
rsync -a \
  --exclude '.git' \
  --exclude 'rpmbuild' \
  --exclude '*.rpm' \
  --exclude 'dist' \
  --exclude 'build' \
  ./ "${PKGROOT}/"

cp packaging/${APP_NAME}.desktop "${PKGROOT}/packaging/${APP_NAME}.desktop"
cp packaging/${APP_NAME}.png "${PKGROOT}/packaging/${APP_NAME}.png"

mkdir -p "${TOPDIR}/SOURCES" "${TOPDIR}/SPECS" "${TOPDIR}/BUILD" "${TOPDIR}/BUILDROOT" "${TOPDIR}/RPMS" "${TOPDIR}/SRPMS"

tar -C "${WORKDIR}" -czf "${TOPDIR}/SOURCES/${APP_NAME}-${VERSION}.tar.gz" "${APP_NAME}-${VERSION}"
cp packaging/${APP_NAME}.spec "${TOPDIR}/SPECS/${APP_NAME}.spec"

rpmbuild -ba "${TOPDIR}/SPECS/${APP_NAME}.spec"

echo
echo "Built packages:"
find "${TOPDIR}/RPMS" "${TOPDIR}/SRPMS" -type f \( -name '*.rpm' -o -name '*.src.rpm' \) | sort
