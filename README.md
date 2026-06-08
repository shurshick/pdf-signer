# PDF Signer

Desktop utility for selecting PDF files, reading certificate metadata from
CryptoPro, generating a visual electronic-signature stamp, placing it into PDF,
and creating a detached `.sig` signature for the stamped PDF.

The application is written in Go with Fyne and targets Linux x86_64
workstations with CryptoPro CSP installed. Windows is useful for editing the
code, but it is not a supported build or runtime target for this project.

Multiple PDF files can be added to the signing queue and processed in one run.
Duplicate files are ignored. Output files use the `_stamped.pdf` suffix. They
can be saved next to each source PDF or into a selected output folder.

The visible stamp is added to all pages before signing. The detached `.sig`
signature is then created for the stamped PDF, so the signature target matches
the document that includes the visible stamp.

The interface, dialogs, errors, and visible stamp are localized automatically:
Russian is used for `ru*` locales, and English is used for all other locales.

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

## Development Notes

Local builds should be run on a Linux system with the Fyne native dependencies
installed. The GitHub Actions workflow uses Ubuntu only as an automated build
environment for the Linux/amd64 binary; final runtime validation should be done
on the target distribution with CryptoPro CSP.

Windows-style embedded PDF signatures are not implemented in this Linux version
yet. The current Linux signing mode creates a detached `.sig` file through
CryptoPro `csptest`.

## Notes

Keep generated binaries, RPM/DEB packages, stamped PDFs, and `.sig` files out of
git; they are ignored by `.gitignore` and should be published as release
artifacts when needed. Published packages are intended for Linux x86_64.

---

# PDF Signer на русском

Настольная утилита для выбора PDF-файлов, чтения данных сертификата из
CryptoPro, создания видимого штампа электронной подписи, добавления этого
штампа в PDF и создания открепленной `.sig` подписи для PDF со штампом.

Приложение написано на Go с использованием Fyne и рассчитано на рабочие станции
Linux x86_64 с установленным CryptoPro CSP. Windows можно использовать для
редактирования кода, но сборка и запуск под Windows в этом проекте не
поддерживаются.

В очередь подписи можно добавить несколько PDF-файлов и обработать их за один
запуск. Дубликаты в очереди игнорируются. Выходные файлы получают суффикс
`_stamped.pdf`. Их можно сохранять рядом с каждым исходным PDF или в выбранную
папку вывода.

Видимый штамп добавляется на все страницы до подписи. После этого открепленная
`.sig` подпись создается для PDF со штампом, поэтому цель подписи совпадает с
документом, который видит пользователь.

Интерфейс, диалоги, сообщения об ошибках и видимый штамп локализуются
автоматически: русский язык включается для locale `ru*`, во всех остальных
случаях используется английский.

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

## Примечания для разработки

Локальную сборку следует выполнять на Linux-системе с установленными нативными
зависимостями Fyne. GitHub Actions использует Ubuntu только как автоматическую
среду сборки Linux/amd64; финальную проверку запуска нужно делать на целевом
дистрибутиве с CryptoPro CSP.

Встроенная PDF-подпись как в Windows-версии пока не реализована в Linux-версии.
Текущий Linux-режим создает открепленный `.sig` файл через CryptoPro `csptest`.

## Важно

Сгенерированные бинарные файлы, RPM/DEB пакеты, PDF со штампом и `.sig` файлы
не нужно хранить в git: они уже добавлены в `.gitignore` и при необходимости
должны публиковаться как release artifacts. Публикуемые пакеты предназначены для
Linux x86_64.
