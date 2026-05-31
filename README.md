# PDF Signer

Desktop utility for selecting a PDF, reading certificate metadata from
CryptoPro, generating a visual electronic-signature stamp, and placing it into a
PDF.

The application is written in Go with Fyne and targets RPM-based Linux
workstations with CryptoPro CSP installed. Windows is useful for editing the
code, but it is not a supported build or runtime target for this project.

## Runtime Requirements

- RPM-based Linux desktop environment
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

Local builds should be run on an RPM-based Linux system with the Fyne native
dependencies installed. The GitHub Actions workflow uses Ubuntu only as an
automated build environment for the Linux/amd64 binary; final runtime validation
should be done on the target RPM-based distribution.

## Notes

The app signs the selected PDF through the native CryptoPro tools and then adds a
visible stamp image to the document. Keep generated binaries, RPM files, and
signed PDFs out of git; they are ignored by `.gitignore` and should be published
as release artifacts when needed. Published RPM packages are intended for
RPM-based Linux distributions on x86_64.

---

# PDF Signer на русском

Настольная утилита для выбора PDF-файла, чтения данных сертификата из CryptoPro,
создания визуального штампа электронной подписи и добавления этого штампа в PDF.

Приложение написано на Go с использованием Fyne и рассчитано на рабочие станции
RPM-based Linux с установленным CryptoPro CSP. Windows можно использовать для
редактирования кода, но сборка и запуск под Windows в этом проекте не
поддерживаются.

## Требования для запуска

- RPM-based Linux с графическим окружением;
- установленный CryptoPro CSP;
- `certmgr` по пути `/opt/cprocsp/bin/amd64/certmgr`;
- `pdfcpu`, доступный через `PATH`.

## Сборка

```bash
go mod download
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o pdfsigner .
```

Запуск:

```bash
./pdfsigner
```

## Сборка RPM

```bash
chmod +x scripts/build-rpm.sh
./scripts/build-rpm.sh
```

Файлы для RPM-пакета находятся в каталоге `packaging/`.

## Примечания для разработки

Локальную сборку следует выполнять на RPM-based Linux системе с установленными
нативными зависимостями Fyne. GitHub Actions использует Ubuntu только как
автоматическую среду сборки Linux/amd64; финальную проверку запуска нужно делать
на целевом RPM-based дистрибутиве.

## Важно

Приложение подписывает выбранный PDF через нативные инструменты CryptoPro, затем
добавляет в документ видимый штамп. Сгенерированные бинарные файлы, RPM-пакеты и
подписанные PDF не нужно хранить в git: они уже добавлены в `.gitignore` и при
необходимости должны публиковаться как release artifacts. Публикуемые RPM-пакеты
предназначены для RPM-based Linux x86_64.
