#!/usr/bin/env bash
set -euo pipefail

APP_NAME="pdfsigner"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILDROOT="${ROOT_DIR}/.rpmbuild"
DIST_DIR="${ROOT_DIR}/dist"

SPEC_TEMPLATE="${ROOT_DIR}/pdfsigner.spec"
SPEC_GENERATED="${BUILDROOT}/SPECS/pdfsigner.spec"

DESKTOP_FILE="${ROOT_DIR}/packaging/pdfsigner.desktop"
ICON_FILE="${ROOT_DIR}/packaging/pdfsigner.png"
BIN_FILE="${ROOT_DIR}/${APP_NAME}"

DEFAULT_VERSION="0.1.2"
DEFAULT_RELEASE="1"
DEFAULT_CHANGELOG="Updated package"

read -r -p "Package version [${DEFAULT_VERSION}]: " VERSION
VERSION="${VERSION:-$DEFAULT_VERSION}"

read -r -p "Release [${DEFAULT_RELEASE}]: " RELEASE
RELEASE="${RELEASE:-$DEFAULT_RELEASE}"

echo "Enter changelog lines for %changelog."
echo "Finish input with an empty line:"
CHANGELOG_LINES=()
while IFS= read -r line; do
  [[ -z "$line" ]] && break
  CHANGELOG_LINES+=("$line")
done

if [[ ${#CHANGELOG_LINES[@]} -eq 0 ]]; then
  CHANGELOG_LINES=("${DEFAULT_CHANGELOG}")
fi

CHANGELOG_TEXT=""
for line in "${CHANGELOG_LINES[@]}"; do
  CHANGELOG_TEXT="${CHANGELOG_TEXT}- ${line}"$'\n'
done
CHANGELOG_TEXT="${CHANGELOG_TEXT%$'\n'}"

CHANGELOG_DATE="$(LC_TIME=C date '+%a %b %d %Y')"
CHANGELOG_USER="${RPM_PACKAGER:-shurshick <noreply@example.com>}"

echo "==> Checking project files"

for required_file in "${SPEC_TEMPLATE}" "${DESKTOP_FILE}" "${ICON_FILE}"; do
  if [[ ! -f "${required_file}" ]]; then
    echo "Error: ${required_file} not found"
    exit 1
  fi
done

echo "==> Cleaning previous build output"
rm -rf "${BUILDROOT}" "${DIST_DIR}"
mkdir -p "${BUILDROOT}"/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
mkdir -p "${DIST_DIR}"

echo "==> Building Go binary"
cd "${ROOT_DIR}"
go mod download
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o "${BIN_FILE}" .

if [[ ! -f "${BIN_FILE}" ]]; then
  echo "Error: ${BIN_FILE} was not built"
  exit 1
fi

echo "==> Preparing RPM sources"
cp -f "${BIN_FILE}" "${BUILDROOT}/SOURCES/"
mkdir -p "${BUILDROOT}/SOURCES/packaging"
cp -f "${DESKTOP_FILE}" "${BUILDROOT}/SOURCES/packaging/"
cp -f "${ICON_FILE}" "${BUILDROOT}/SOURCES/packaging/"

echo "==> Generating SPEC"
awk -v version="${VERSION}" -v release="${RELEASE}" '
{
  if ($1 == "Version:") {
    print "Version:        " version
  } else if ($1 == "Release:") {
    print "Release:        " release
  } else {
    print
  }
}
' "${SPEC_TEMPLATE}" > "${SPEC_GENERATED}"

cat >> "${SPEC_GENERATED}" <<EOF

%changelog
* ${CHANGELOG_DATE} ${CHANGELOG_USER} ${VERSION}-${RELEASE}
${CHANGELOG_TEXT}
EOF

echo "==> Building RPM"
rpmbuild -bb "${SPEC_GENERATED}" \
  --define "_topdir ${BUILDROOT}" \
  --define "_sourcedir ${BUILDROOT}/SOURCES" \
  --define "_specdir ${BUILDROOT}/SPECS" \
  --define "_builddir ${BUILDROOT}/BUILD" \
  --define "_buildrootdir ${BUILDROOT}/BUILDROOT" \
  --define "_rpmdir ${BUILDROOT}/RPMS" \
  --define "_srcrpmdir ${BUILDROOT}/SRPMS"

echo "==> Copying RPM to dist/"
find "${BUILDROOT}/RPMS" -type f -name "*.rpm" -exec cp -f {} "${DIST_DIR}/" \;

echo "==> Done"
echo "Version: ${VERSION}"
echo "Release: ${RELEASE}"
echo "RPM packages:"
ls -lh "${DIST_DIR}"
