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

DEFAULT_VERSION="0.1.0"
DEFAULT_RELEASE="1.red80"
DEFAULT_CHANGELOG="Updated package"

read -r -p "Версия пакета [${DEFAULT_VERSION}]: " VERSION
VERSION="${VERSION:-$DEFAULT_VERSION}"

read -r -p "Release [${DEFAULT_RELEASE}]: " RELEASE
RELEASE="${RELEASE:-$DEFAULT_RELEASE}"

echo "Введите описание изменений для %changelog."
echo "Завершите ввод пустой строкой:"
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
CHANGELOG_USER="${RPM_PACKAGER:-OpenAI <noreply@example.com>}"

echo "==> Проверка файлов проекта"

if [[ ! -f "${SPEC_TEMPLATE}" ]]; then
  echo "Ошибка: не найден ${SPEC_TEMPLATE}"
  exit 1
fi

if [[ ! -f "${DESKTOP_FILE}" ]]; then
  echo "Ошибка: не найден ${DESKTOP_FILE}"
  exit 1
fi

if [[ ! -f "${ICON_FILE}" ]]; then
  echo "Ошибка: не найден ${ICON_FILE}"
  exit 1
fi

echo "==> Очистка старых каталогов"
rm -rf "${BUILDROOT}" "${DIST_DIR}"
mkdir -p "${BUILDROOT}"/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
mkdir -p "${DIST_DIR}"

echo "==> Сборка Go-бинарника"
cd "${ROOT_DIR}"
go mod tidy
go build -o "${BIN_FILE}"

if [[ ! -f "${BIN_FILE}" ]]; then
  echo "Ошибка: бинарник ${BIN_FILE} не собран"
  exit 1
fi

echo "==> Подготовка SOURCES"
cp -f "${BIN_FILE}" "${BUILDROOT}/SOURCES/"
mkdir -p "${BUILDROOT}/SOURCES/packaging"
cp -f "${DESKTOP_FILE}" "${BUILDROOT}/SOURCES/packaging/"
cp -f "${ICON_FILE}" "${BUILDROOT}/SOURCES/packaging/"

echo "==> Генерация SPEC"
awk -v version="${VERSION}" -v release="${RELEASE}" '
BEGIN { doneVersion=0; doneRelease=0 }
{
  if ($1 == "Version:") {
    print "Version:        " version
    doneVersion=1
  } else if ($1 == "Release:") {
    print "Release:        " release
    doneRelease=1
  } else {
    print
  }
}
END {
  if (!doneVersion) print "Version:        " version > "/dev/stderr"
  if (!doneRelease) print "Release:        " release > "/dev/stderr"
}
' "${SPEC_TEMPLATE}" > "${SPEC_GENERATED}"

cat >> "${SPEC_GENERATED}" <<EOF

%changelog
* ${CHANGELOG_DATE} ${CHANGELOG_USER} ${VERSION}-${RELEASE}
${CHANGELOG_TEXT}
EOF

echo "==> Сборка RPM"
rpmbuild -bb "${SPEC_GENERATED}" \
  --define "_topdir ${BUILDROOT}" \
  --define "_sourcedir ${BUILDROOT}/SOURCES" \
  --define "_specdir ${BUILDROOT}/SPECS" \
  --define "_builddir ${BUILDROOT}/BUILD" \
  --define "_buildrootdir ${BUILDROOT}/BUILDROOT" \
  --define "_rpmdir ${BUILDROOT}/RPMS" \
  --define "_srcrpmdir ${BUILDROOT}/SRPMS"

echo "==> Копирование RPM в dist/"
find "${BUILDROOT}/RPMS" -type f -name "*.rpm" -exec cp -f {} "${DIST_DIR}/" \;

echo "==> Готово"
echo "Версия: ${VERSION}"
echo "Release: ${RELEASE}"
echo "RPM пакеты:"
ls -lh "${DIST_DIR}"