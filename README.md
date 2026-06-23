# PDF Signer на русском

Настольная утилита для выбора PDF-файлов, чтения данных сертификата из
CryptoPro, создания видимого штампа электронной подписи, добавления этого
штампа в PDF и создания цифровой подписи — встроенной в PDF или открепленной
`.sig` файлом.

Приложение написано на Go с использованием Fyne и рассчитано на рабочие станции
Linux x86_64 с установленным CryptoPro CSP.

## Возможности

- выбор одного или нескольких PDF-файлов;
- выбор действующего сертификата из хранилища CryptoPro (`uMy`);
- видимый штамп подписи на всех страницах PDF;
- отпечаток сертификата на штампе;
- **встроенная PDF-подпись** — CAdES-BES подпись встраивается внутрь PDF;
- **открепленная `.sig` подпись** — подпись сохраняется отдельным файлом;
- **режим «оба»** — одновременно встроенная подпись и открепленный `.sig` файл;
- **редактор штампа** — настройка шаблона, позиции, размера, шрифта, полей;
- **проверка подписи** — верификация `.sig` файлов и встроенных подписей;
- **диагностика CryptoPro** — проверка certmgr, csptest, сертификатов;
- **дружественные сообщения об ошибках** — понятные описание вместо сырых исключений;
- **логирование** — запись действий в файл с санитизацией секретов;
- **настройки** — сохранение профиля штампа, экспорт/импорт настроек JSON;
- **автоматическое размещение штампа** — избегает перекрытия текста на странице;
- **логотип в штампе** — PNG/JPG логотип в углу штампа;
- сохранение результата в указанную папку или рядом с исходным PDF;
- окно «О приложении» с версией, копирайтом, ссылкой на проект и проверкой обновлений;
- пакетное подписание нескольких файлов за один запуск;
- автоматический русский или английский интерфейс по языку системы;
- сборка RPM и DEB через GitHub Actions.

## Режимы подписания

| Режим | Описание |
|---|---|
| **Встроенная PDF-подпись** | CAdES-BES подпись встраивается внутрь PDF-документа через `csptest`. |
| **Открепленный `.sig`** | Подпись сохраняется отдельным `.sig` файлом рядом с PDF. |
| **Оба** | Создается и встроенная подпись в PDF, и отдельный `.sig` файл. |

## Диагностика

Окно диагностики проверяет:
- наличие `certmgr` и `csptest` по стандартным путям;
- наличие сертификатов в хранилище `uMy`;
- количество сертификатов с закрытыми ключами;
- готовность к подписанию.

## Редактор штампа

Три встроенных профиля:
- **Минимальный** — 70×25 мм, основные поля;
- **Стандартный** — 90×35 мм, все основные поля;
- **Подробный** — 120×45 мм, все поля включая ИНН, СНИЛС, отпечаток.

Возможности:
- настройка шаблона, страниц, позиции, размера;
- управление видимостью полей (владелец, издатель, дата, причина, серийный номер);
- шрифт и непрозрачность;
- автоматическое размещение штампа (избегает текста);
- логотип в штампе (PNG/JPG, масштаб 100–300%);
- загрузка и сохранение профилей в JSON.

## Проверка подписи

- проверка `.sig` файлов (поиск соответствующего PDF);
- проверка встроенных подписей в PDF через `csptest -sfsign -verify`;
- подробный отчёт с уровнем VALID/WARNING/INVALID;
- экспорт отчёта в TXT.

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

## Сборка RPM

```bash
chmod +x build-rpm.sh
./build-rpm.sh
```

## Сборка DEB

```bash
chmod +x build-deb.sh
./build-deb.sh
```

## Тесты

```bash
go test ./...
```

## Лицензия

AGPL-3.0-or-later.

---

# PDF Signer

Desktop utility for selecting PDF files, reading certificate metadata from
CryptoPro, generating a visual electronic-signature stamp, placing it into PDF,
and creating a digital signature — either embedded in the PDF or as a detached
`.sig` file.

The application is written in Go with Fyne and targets Linux x86_64
workstations with CryptoPro CSP installed.

## Features

- select one or more PDF files;
- select a currently valid certificate from the CryptoPro `uMy` store;
- visible signature stamp on every page;
- certificate thumbprint on the stamp;
- **embedded PDF signature** — CAdES-BES signature embedded inside the PDF;
- **detached `.sig` signature** — signature saved as a separate `.sig` file;
- **both mode** — embedded signature plus a separate `.sig` file;
- **stamp editor** — configure template, position, size, font, fields;
- **signature verification** — verify `.sig` files and embedded signatures;
- **CryptoPro diagnostics** — check certmgr, csptest, certificate store;
- **friendly error messages** — human-readable error descriptions;
- **application logging** — timestamped file logging with secret sanitization;
- **settings persistence** — stamp profile saved, JSON export/import;
- **smart stamp placement** — avoids overlapping text on PDF pages;
- **logo support** — PNG/JPG logo in stamp corner;
- save outputs to selected folder or next to source PDF;
- About dialog with version, copyright, project link, and update check;
- batch-sign multiple PDFs in one run;
- Russian or English UI from system language;
- RPM and DEB packages via GitHub Actions.

## Signing Modes

| Mode | Description |
|---|---|
| **Embedded PDF signature** | CAdES-BES signature embedded inside the PDF via `csptest`. |
| **Detached `.sig`** | Signature saved as a separate `.sig` file alongside the PDF. |
| **Both** | Both an embedded signature and a separate `.sig` file. |

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

## Tests

```bash
go test ./...
```

## License

AGPL-3.0-or-later.
