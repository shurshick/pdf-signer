# PDF Signer

Desktop utility for selecting a PDF, reading certificate metadata from
CryptoPro, generating a visual electronic-signature stamp, and placing it into a
PDF.

The application is written in Go with Fyne and targets Linux workstations with
CryptoPro CSP installed. Windows is useful for editing the code, but it is not a
supported build or runtime target for this project.

## Runtime Requirements

- Linux desktop environment
- CryptoPro CSP
- `certmgr` at `/opt/cprocsp/bin/amd64/certmgr`
- `pdfcpu` available in `PATH`

## Build

```bash
go mod download
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o pdfsigner .
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

## Development Notes

Local builds should be run on Linux with the Fyne native dependencies installed.
The GitHub Actions workflow uses Ubuntu and validates the intended Linux/amd64
target.

## Notes

The app signs the selected PDF through the native CryptoPro tools and then adds a
visible stamp image to the document. Keep generated binaries, RPM files, and
signed PDFs out of git; they are ignored by `.gitignore` and should be published
as release artifacts when needed.
