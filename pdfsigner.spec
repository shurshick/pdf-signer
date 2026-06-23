Name:           pdfsigner
Version:        0.4.0
Release:        1
Summary:        PDF signer and visual electronic signature stamp tool

License:        AGPLv3
URL:            https://github.com/shurshick/pdf-signer
BuildArch:      x86_64

Requires:       /opt/cprocsp/bin/amd64/certmgr
Requires:       /opt/cprocsp/bin/amd64/csptest

%description
Desktop application for adding a visual electronic signature stamp to PDF
documents and creating embedded or detached digital signatures using
CryptoPro CSP. Supports embedded PDF signatures, detached .sig files, or both.

%prep
# Ничего не распаковываем.

%build
# Сборка бинарника выполняется заранее внешним скриптом.

%install
rm -rf %{buildroot}

install -D -m 0755 %{_sourcedir}/pdfsigner %{buildroot}/usr/bin/pdfsigner
install -D -m 0644 %{_sourcedir}/packaging/pdfsigner.desktop %{buildroot}/usr/share/applications/pdfsigner.desktop
install -D -m 0644 %{_sourcedir}/packaging/pdfsigner.png %{buildroot}/usr/share/icons/hicolor/256x256/apps/pdfsigner.png

%files
/usr/bin/pdfsigner
/usr/share/applications/pdfsigner.desktop
/usr/share/icons/hicolor/256x256/apps/pdfsigner.png
