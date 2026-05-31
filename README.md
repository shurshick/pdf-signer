# PDF Signer

Desktop utility for selecting a PDF, reading certificate metadata from
CryptoPro, generating a visual electronic-signature stamp, and placing it into a
PDF.

The application is written in Go with Fyne and is currently focused on Linux
workstations with CryptoPro CSP installed.

## Runtime Requirements

- Linux desktop environment
- CryptoPro CSP
- `certmgr` at `/opt/cprocsp/bin/amd64/certmgr`
- `pdfcpu` available in `PATH`

## Build

```bash
go mod download
go build -o pdfsigner .
```

Run:

```bash
./pdfsigner
```

## RPM Build

```bash
chmod +x scripts/build-rpm.sh
./scripts/build-rpm.sh
```

The RPM packaging assets live in `packaging/`.

## Notes

The app signs the selected PDF through the native CryptoPro tools and then adds a
visible stamp image to the document. Keep generated binaries, RPM files, and
signed PDFs out of git; they are ignored by `.gitignore` and should be published
as release artifacts when needed.
