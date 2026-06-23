# PDF Signer на русском

Настольная утилита для выбора PDF-файлов, чтения данных сертификата из
CryptoPro, создания видимого штампа электронной подписи, добавления этого
штампа в PDF и создания цифровой подписи — встроенной в PDF или открепленной
`.sig` файлом.

Приложение написано на Go с использованием Fyne и рассчитано на рабочие станции
Linux x86_64 с установленным CryptoPro CSP. Windows можно использовать для
редактирования кода, но сборка и запуск под Windows в этом проекте не
поддерживаются.

## Возможности

- выбор одного или нескольких PDF-файлов;
- выбор действующего сертификата из хранилища CryptoPro (`uMy`);
- видимый штамп подписи на всех страницах PDF;
- отпечаток сертификата на штампе;
- **встроенная PDF-подпись** — CAdES-BES подпись встраивается внутрь PDF;
- **открепленная `.sig` подпись** — подпись сохраняется отдельным файлом;
- **режим «оба»** — одновременно встроенная подпись и открепленный `.sig` файл;
- сохранение результата в указанную папку или рядом с исходным PDF;
- окно «О приложении» с версией, копирайтом, ссылкой на проект и проверкой обновлений;
- иконка приложения;
- пакетное подписание нескольких файлов за один запуск;
- автоматический русский или английский интерфейс по языку системы;
- сборка RPM и DEB через GitHub Actions.

## Режимы подписания

| Режим | Описание |
|---|---|
| **Встроенная PDF-подпись** | CAdES-BES подпись встраивается внутрь PDF-документа через `csptest`. PDF получает встроенную цифровую подпись. |
| **Открепленный `.sig`** | Подпись сохраняется отдельным `.sig` файлом рядом с PDF. PDF содержит только видимый штамп. |
| **Оба** | Создается и встроенная подпись в PDF, и отдельный `.sig` файл. |

## Требования для запуска

- Linux x86_64 с графическим окружением
- установленный CryptoPro CSP
- `certmgr` по пути `/opt/cprocsp/bin/amd64/certmgr`
- `csptest` по пути `/opt/cprocsp/bin/amd64/csptest`
- `pdfcpu`, доступный через `PATH`

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
chmod +x build-rpm.sh
./build-rpm.sh
```

Файлы для RPM-пакета находятся в каталоге `packaging/`.

## Сборка DEB

```bash
chmod +x build-deb.sh
./build-deb.sh
```

## Тесты

```bash
go test ./...
```

## Примечания для разработки

Локальную сборку следует выполнять на Linux-системе с установленными нативными
зависимостями Fyne. GitHub Actions использует Ubuntu только как автоматическую
среду сборки Linux/amd64; финальную проверку запуска нужно делать на целевом
дистрибутиве с CryptoPro CSP.

CryptoPro CSP, сертификаты и ключевые контейнеры не входят в поставку приложения.

## Лицензия

AGPL-3.0-or-later.

## Важно

Сгенерированные бинарные файлы, RPM/DEB пакеты, PDF со штампом и `.sig` файлы
не нужно хранить в git: они уже добавлены в `.gitignore` и при необходимости
должны публиковаться как release artifacts. Публикуемые пакеты предназначены для
Linux x86_64.

---

# PDF Signer

Desktop utility for selecting PDF files, reading certificate metadata from
CryptoPro, generating a visual electronic-signature stamp, placing it into PDF,
and creating a digital signature — either embedded in the PDF or as a detached
`.sig` file.

The application is written in Go with Fyne and targets Linux x86_64
workstations with CryptoPro CSP installed. Windows is useful for editing the
code, but it is not a supported build or runtime target for this project.

## Features

- select one or more PDF files;
- select a currently valid certificate from the CryptoPro `uMy` store;
- visible signature stamp on every page;
- certificate thumbprint on the stamp;
- **embedded PDF signature** — CAdES-BES signature embedded inside the PDF;
- **detached `.sig` signature** — signature saved as a separate `.sig` file;
- **both mode** — embedded signature plus a separate `.sig` file;
- save outputs either to the selected output folder or next to each source PDF;
- About dialog with version, copyright, project link, and update check;
- application icon;
- batch-sign multiple PDFs in one run;
- choose Russian or English UI automatically from the system UI language;
- build RPM and DEB packages with GitHub Actions.

## Signing Modes

| Mode | Description |
|---|---|
| **Embedded PDF signature** | CAdES-BES signature embedded inside the PDF document via `csptest`. The PDF gets an embedded digital signature. |
| **Detached `.sig`** | Signature saved as a separate `.sig` file alongside the PDF. The PDF contains only the visible stamp. |
| **Both** | Both an embedded signature in the PDF and a separate `.sig` file are created. |

## Runtime Requirements

- Linux x86_64 desktop environment
- CryptoPro CSP
- `certmgr` at `/opt/cprocsp/bin/amd64/certmgr`
- `csptest` at `/opt/cprocsp/bin/amd64/csptest`
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
chmod +x build-rpm.sh
./build-rpm.sh
```

The RPM packaging assets live in `packaging/`.

## DEB Build

```bash
chmod +x build-deb.sh
./build-deb.sh
```

## Tests

```bash
go test ./...
```

## Development Notes

Local builds should be run on a Linux system with the Fyne native dependencies
installed. The GitHub Actions workflow uses Ubuntu only as an automated build
environment for the Linux/amd64 binary; final runtime validation should be done
on the target distribution with CryptoPro CSP.

CryptoPro CSP, certificates, and key containers are not bundled with this application.

## License

AGPL-3.0-or-later.

## Notes

Keep generated binaries, RPM/DEB packages, stamped PDFs, and `.sig` files out of
git; they are ignored by `.gitignore` and should be published as release
artifacts when needed. Published packages are intended for Linux x86_64.
