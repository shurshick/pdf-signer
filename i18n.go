package main

import (
	"os"
	"strings"
)

type messageID string

const (
	msgWindowTitle            messageID = "window_title"
	msgPDFNotSelected         messageID = "pdf_not_selected"
	msgOutPDFNotSelected      messageID = "out_pdf_not_selected"
	msgCertNotSelected        messageID = "cert_not_selected"
	msgOwner                  messageID = "owner"
	msgIssuer                 messageID = "issuer"
	msgSelectPDF              messageID = "select_pdf"
	msgClearPDFs              messageID = "clear_pdfs"
	msgSelectOutputFolder     messageID = "select_output_folder"
	msgBrowse                 messageID = "browse"
	msgAbout                  messageID = "about"
	msgVersion                messageID = "version"
	msgProjectLink            messageID = "project_link"
	msgCheckUpdates           messageID = "check_updates"
	msgUpdateAvailableTitle   messageID = "update_available_title"
	msgUpdateAvailableBody    messageID = "update_available_body"
	msgNoUpdatesTitle         messageID = "no_updates_title"
	msgNoUpdatesBody          messageID = "no_updates_body"
	msgUpdateCheckFailed      messageID = "update_check_failed"
	msgCurrentVersion         messageID = "current_version"
	msgLatestVersion          messageID = "latest_version"
	msgDownload               messageID = "download"
	msgClose                  messageID = "close"
	msgReason                 messageID = "reason"
	msgDefaultReason          messageID = "default_reason"
	msgSaveNextToSource       messageID = "save_next_to_source"
	msgDetachedModeNote       messageID = "detached_mode_note"
	msgSignAndStamp           messageID = "sign_and_stamp"
	msgError                  messageID = "error"
	msgChoosePDF              messageID = "choose_pdf"
	msgChooseCert             messageID = "choose_cert"
	msgChooseOutPDF           messageID = "choose_out_pdf"
	msgDone                   messageID = "done"
	msgSignature              messageID = "signature"
	msgStampedPDF             messageID = "stamped_pdf"
	msgSelectedPDFs           messageID = "selected_pdfs"
	msgProcessedFiles         messageID = "processed_files"
	msgBatchOutputNote        messageID = "batch_output_note"
	msgCertificate            messageID = "certificate"
	msgScale                  messageID = "scale"
	msgCertmgrError           messageID = "certmgr_error"
	msgNoCerts                messageID = "no_certs"
	msgNoPDF                  messageID = "no_pdf"
	msgAbsPDFError            messageID = "abs_pdf_error"
	msgAbsSignatureError      messageID = "abs_signature_error"
	msgEmptyCertCN            messageID = "empty_cert_cn"
	msgSignError              messageID = "sign_error"
	msgSignatureMissing       messageID = "signature_missing"
	msgOutputSignatureMissing messageID = "output_signature_missing"
	msgInputPDFMissing        messageID = "input_pdf_missing"
	msgOutputPDFMissing       messageID = "output_pdf_missing"
	msgStampFileMissing       messageID = "stamp_file_missing"
	msgStampPDFError          messageID = "stamp_pdf_error"
	msgStampTitle             messageID = "stamp_title"
	msgDate                   messageID = "date"
	msgSerialNumber           messageID = "serial_number"
	msgSignatureShort         messageID = "signature_short"
	msgSigningMode            messageID = "signing_mode"
	msgModeEmbedded           messageID = "mode_embedded"
	msgModeDetached           messageID = "mode_detached"
	msgModeBoth               messageID = "mode_both"
	msgModeEmbeddedDesc       messageID = "mode_embedded_desc"
	msgModeDetachedDesc       messageID = "mode_detached_desc"
	msgModeBothDesc           messageID = "mode_both_desc"
	msgSignEmbeddedError      messageID = "sign_embedded_error"
	msgSignedPDF              messageID = "signed_pdf"
	msgEmbeddedSignature      messageID = "embedded_signature"
	msgDiagnostics            messageID = "diagnostics"
	msgDiagnosticsTitle       messageID = "diagnostics_title"
	msgRunDiagnostics         messageID = "run_diagnostics"
	msgCopyReport             messageID = "copy_report"
	msgSaveReport             messageID = "save_report"
	msgOpenLogsFolder         messageID = "open_logs_folder"
	msgDiagnosticsRunning     messageID = "diagnostics_running"
	msgCertmgrNotFound        messageID = "certmgr_not_found"
	msgCsptestNotFound        messageID = "csptest_not_found"
	msgNoCertsFound           messageID = "no_certs_found"
	msgCertCNEmpty            messageID = "cert_cn_empty"
	msgPDFNotSpecified        messageID = "pdf_not_specified"
	msgSigningError           messageID = "signing_error"
	msgEmbeddedSignError      messageID = "embedded_sign_error"
	msgStampError             messageID = "stamp_error"
	msgPermissionDenied       messageID = "permission_denied"
	msgCertExpired            messageID = "cert_expired"
	msgCertNotValidYet        messageID = "cert_not_valid_yet"
	msgPrivateKeyMissing      messageID = "private_key_missing"
	msgCryptoError            messageID = "crypto_error"
	msgVerifySignature        messageID = "verify_signature"
	msgVerificationTitle      messageID = "verification_title"
	msgRunVerification        messageID = "run_verification"
	msgExportTxt              messageID = "export_txt"
	msgVerifyAfterSigning     messageID = "verify_after_signing"
	msgStampEditor            messageID = "stamp_editor"
	msgStampEditorTitle       messageID = "stamp_editor_title"
	msgTemplateMinimal        messageID = "template_minimal"
	msgTemplateStandard       messageID = "template_standard"
	msgTemplateDetailed       messageID = "template_detailed"
	msgStampTemplate          messageID = "stamp_template"
	msgStampPages             messageID = "stamp_pages"
	msgStampPosition          messageID = "stamp_position"
	msgStampSize              messageID = "stamp_size"
	msgWidthMm                messageID = "width_mm"
	msgHeightMm               messageID = "height_mm"
	msgFontSize               messageID = "font_size"
	msgMinFontSize            messageID = "min_font_size"
	msgOpacity                messageID = "opacity"
	msgPosBottomRight         messageID = "pos_bottom_right"
	msgPosBottomLeft          messageID = "pos_bottom_left"
	msgPosTopRight            messageID = "pos_top_right"
	msgPosTopLeft             messageID = "pos_top_left"
	msgLoadProfile            messageID = "load_profile"
	msgSaveProfile            messageID = "save_profile"
	msgStampTooSmall          messageID = "stamp_too_small"
	msgFontSizeTooSmall       messageID = "font_size_too_small"
	msgOpacityLow             messageID = "opacity_low"
	msgAutoPlaceStamp         messageID = "auto_place_stamp"
	msgLogoPath               messageID = "logo_path"
	msgLogoScale              messageID = "logo_scale"
	msgChooseLogo             messageID = "choose_logo"
	msgRemoveLogo             messageID = "remove_logo"
	msgLogoMissing            messageID = "logo_missing"
	msgLogoTooLarge           messageID = "logo_too_large"
	msgExportSettings         messageID = "export_settings"
	msgImportSettings         messageID = "import_settings"
	msgSettingsExported       messageID = "settings_exported"
	msgSettingsImported       messageID = "settings_imported"
	msgFileSummary            messageID = "file_summary"
	msgVerificationReport     messageID = "verification_report"
	msgSignatureValid         messageID = "signature_valid"
	msgSignatureInvalid       messageID = "signature_invalid"
	msgPreview                messageID = "preview"
	msgReady                  messageID = "ready"
)

var messages = map[messageID][2]string{
	msgWindowTitle:            {"PDF signer + stamp", "PDF подпись + штамп"},
	msgPDFNotSelected:         {"PDF is not selected", "PDF не выбран"},
	msgOutPDFNotSelected:      {"Output folder is not selected", "Папка вывода не выбрана"},
	msgCertNotSelected:        {"Certificate is not selected", "Сертификат не выбран"},
	msgOwner:                  {"Owner", "Владелец"},
	msgIssuer:                 {"Issuer", "Издатель"},
	msgSelectPDF:              {"Select PDF", "Выбрать PDF"},
	msgClearPDFs:              {"Clear PDFs", "Очистить PDF"},
	msgSelectOutputFolder:     {"Output folder", "Папка вывода"},
	msgBrowse:                 {"Browse", "Обзор"},
	msgAbout:                  {"About", "О приложении"},
	msgVersion:                {"Version", "Версия"},
	msgProjectLink:            {"Project", "Проект"},
	msgCheckUpdates:           {"Check for updates", "Проверить обновление"},
	msgUpdateAvailableTitle:   {"Update available", "Доступно обновление"},
	msgUpdateAvailableBody:    {"A newer version is available. Open the release page to download it?", "Доступна новая версия. Открыть страницу релиза для скачивания?"},
	msgNoUpdatesTitle:         {"No updates", "Обновлений нет"},
	msgNoUpdatesBody:          {"You are using the latest available version.", "У вас установлена последняя доступная версия."},
	msgUpdateCheckFailed:      {"Failed to check for updates", "Не удалось проверить обновления"},
	msgCurrentVersion:         {"Current version", "Текущая версия"},
	msgLatestVersion:          {"Latest version", "Последняя версия"},
	msgDownload:               {"Download", "Скачать"},
	msgClose:                  {"Close", "Закрыть"},
	msgReason:                 {"Reason", "Причина"},
	msgDefaultReason:          {"Document signing", "Подписание документа"},
	msgSaveNextToSource:       {"Save next to source PDF", "Сохранять рядом с исходным PDF"},
	msgDetachedModeNote:       {"Creates a detached .sig for the stamped PDF.", "Создает открепленную .sig для PDF со штампом."},
	msgSignAndStamp:           {"Sign and stamp", "Подписать и поставить штамп"},
	msgError:                  {"Error", "Ошибка"},
	msgChoosePDF:              {"Select a PDF file", "Выберите PDF"},
	msgChooseCert:             {"Select a certificate", "Выберите сертификат"},
	msgChooseOutPDF:           {"Choose output folder", "Укажите папку вывода"},
	msgDone:                   {"Done", "Готово"},
	msgSignature:              {"Signature", "Подпись"},
	msgStampedPDF:             {"Stamped PDF", "PDF со штампом"},
	msgSelectedPDFs:           {"Selected PDFs", "Выбранные PDF"},
	msgProcessedFiles:         {"Processed files", "Обработано файлов"},
	msgBatchOutputNote:        {"Output files get a suffix. In detached mode, a .sig is also created.", "Выходные файлы получают суффикс. В режиме .sig также создается открепленная подпись."},
	msgCertificate:            {"Certificate", "Сертификат"},
	msgScale:                  {"Scale", "Масштаб"},
	msgCertmgrError:           {"certmgr error", "ошибка certmgr"},
	msgNoCerts:                {"no certificates found in uMy store", "сертификаты в хранилище uMy не найдены"},
	msgNoPDF:                  {"PDF is not specified", "не указан PDF"},
	msgAbsPDFError:            {"failed to resolve absolute PDF path", "не удалось получить абсолютный путь к PDF"},
	msgAbsSignatureError:      {"failed to resolve absolute signature path", "не удалось получить абсолютный путь к подписи"},
	msgEmptyCertCN:            {"certificate CN is empty", "у сертификата пустой CN"},
	msgSignError:              {"signing error", "ошибка подписи"},
	msgSignatureMissing:       {"signature file was not created", "файл подписи не создан"},
	msgOutputSignatureMissing: {"output signature path is not specified", "не указан выходной файл подписи"},
	msgInputPDFMissing:        {"input PDF is not specified", "не указан входной PDF"},
	msgOutputPDFMissing:       {"output PDF is not specified", "не указан выходной PDF"},
	msgStampFileMissing:       {"stamp file is not specified", "не указан файл штампа"},
	msgStampPDFError:          {"failed to add stamp to PDF", "ошибка добавления штампа в PDF"},
	msgStampTitle:             {"ELECTRONIC SIGNATURE", "ЭЛЕКТРОННАЯ ПОДПИСЬ"},
	msgDate:                   {"Date", "Дата"},
	msgSerialNumber:           {"Serial number", "Серийный номер"},
	msgSignatureShort:         {"ES", "ЭП"},
	msgSigningMode:            {"Signing mode", "Режим подписания"},
	msgModeEmbedded:           {"Embedded PDF signature", "Встроенная PDF-подпись"},
	msgModeDetached:           {"Detached .sig file", "Открепленный .sig файл"},
	msgModeBoth:               {"Both (embedded + .sig)", "Оба (встроенная + .sig)"},
	msgModeEmbeddedDesc:       {"Signature is embedded inside the PDF document.", "Подпись встроена внутрь PDF-документа."},
	msgModeDetachedDesc:       {"Signature is saved as a separate .sig file alongside the PDF.", "Подпись сохраняется отдельным .sig файлом рядом с PDF."},
	msgModeBothDesc:           {"Embedded signature in the PDF plus a separate .sig file.", "Встроенная подпись в PDF и отдельный .sig файл."},
	msgSignEmbeddedError:      {"embedded signing error", "ошибка встроенного подписания"},
	msgSignedPDF:              {"Signed PDF", "Подписанный PDF"},
	msgEmbeddedSignature:      {"Embedded signature", "Встроенная подпись"},
	msgDiagnostics:            {"Diagnostics", "Диагностика"},
	msgDiagnosticsTitle:       {"CryptoPro Diagnostics", "Диагностика CryptoPro"},
	msgRunDiagnostics:         {"Run diagnostics", "Запустить диагностику"},
	msgCopyReport:             {"Copy report", "Копировать отчёт"},
	msgSaveReport:             {"Save report", "Сохранить отчёт"},
	msgOpenLogsFolder:         {"Open logs folder", "Открыть папку логов"},
	msgDiagnosticsRunning:     {"Diagnostics are running...", "Диагностика выполняется..."},
	msgCertmgrNotFound:        {"certmgr not found at /opt/cprocsp/bin/amd64/certmgr", "certmgr не найден по пути /opt/cprocsp/bin/amd64/certmgr"},
	msgCsptestNotFound:        {"csptest not found at /opt/cprocsp/bin/amd64/csptest", "csptest не найден по пути /opt/cprocsp/bin/amd64/csptest"},
	msgNoCertsFound:           {"No certificates found in uMy store", "Сертификаты в хранилище uMy не найдены"},
	msgCertCNEmpty:            {"Certificate CN is empty", "У сертификата пустой CN"},
	msgPDFNotSpecified:        {"PDF is not specified", "Не указан PDF"},
	msgSigningError:           {"Signing error", "Ошибка подписи"},
	msgEmbeddedSignError:      {"Embedded signing error", "Ошибка встроенного подписания"},
	msgStampError:             {"Stamp error", "Ошибка штампа"},
	msgPermissionDenied:       {"Permission denied. Check file permissions.", "Нет прав доступа. Проверьте права на файлы."},
	msgCertExpired:            {"Certificate has expired", "Сертификат истёк"},
	msgCertNotValidYet:        {"Certificate is not yet valid", "Сертификат ещё не действителен"},
	msgPrivateKeyMissing:      {"Private key is missing", "Отсутствует закрытый ключ"},
	msgCryptoError:            {"Cryptographic error", "Криптографическая ошибка"},
	msgVerifySignature:        {"Verify signature", "Проверить подпись"},
	msgVerificationTitle:      {"Signature Verification", "Проверка подписи"},
	msgRunVerification:        {"Verify", "Проверить"},
	msgExportTxt:              {"Export TXT", "Экспорт TXT"},
	msgVerifyAfterSigning:     {"Verify after signing", "Проверять после подписания"},
	msgStampEditor:            {"Stamp editor", "Редактор штампа"},
	msgStampEditorTitle:       {"Stamp Editor", "Редактор штампа"},
	msgTemplateMinimal:        {"Minimal", "Минимальный"},
	msgTemplateStandard:       {"Standard", "Стандартный"},
	msgTemplateDetailed:       {"Detailed", "Подробный"},
	msgStampTemplate:          {"Template", "Шаблон"},
	msgStampPages:             {"Pages", "Страницы"},
	msgStampPosition:          {"Position", "Позиция"},
	msgStampSize:              {"Size", "Размер"},
	msgWidthMm:                {"Width, mm", "Ширина, мм"},
	msgHeightMm:               {"Height, mm", "Высота, мм"},
	msgFontSize:               {"Font size", "Размер шрифта"},
	msgMinFontSize:            {"Min font size", "Мин. размер шрифта"},
	msgOpacity:                {"Opacity, %", "Непрозрачность, %"},
	msgPosBottomRight:         {"Bottom right", "Правый нижний"},
	msgPosBottomLeft:          {"Bottom left", "Левый нижний"},
	msgPosTopRight:            {"Top right", "Правый верхний"},
	msgPosTopLeft:             {"Top left", "Левый верхний"},
	msgLoadProfile:            {"Load profile", "Загрузить профиль"},
	msgSaveProfile:            {"Save profile", "Сохранить профиль"},
	msgStampTooSmall:          {"Stamp may be too small for text", "Штамп может быть слишком мал для текста"},
	msgFontSizeTooSmall:       {"Font size is below minimum", "Размер шрифта ниже минимума"},
	msgOpacityLow:             {"Low opacity may affect readability", "Низкая непрозрачность может повлиять на читаемость"},
	msgAutoPlaceStamp:         {"Auto-place stamp (avoid text)", "Авто-размещение штампа (без перекрытия текста)"},
	msgLogoPath:               {"Logo path", "Путь к логотипу"},
	msgLogoScale:              {"Logo scale, %", "Масштаб логотипа, %"},
	msgChooseLogo:             {"Choose logo", "Выбрать логотип"},
	msgRemoveLogo:             {"Remove logo", "Удалить логотип"},
	msgLogoMissing:            {"Logo file not found", "Файл логотипа не найден"},
	msgLogoTooLarge:           {"Logo file is too large (>1MB)", "Файл логотипа слишком большой (>1MB)"},
	msgExportSettings:         {"Export settings", "Экспорт настроек"},
	msgImportSettings:         {"Import settings", "Импорт настроек"},
	msgSettingsExported:       {"Settings exported.", "Настройки экспортированы."},
	msgSettingsImported:       {"Settings imported.", "Настройки импортированы."},
	msgFileSummary:            {"Selected PDFs: %d", "Выбрано PDF: %d"},
	msgVerificationReport:     {"Verification Report", "Отчёт проверки"},
	msgSignatureValid:         {"Signature is valid", "Подпись действительна"},
	msgSignatureInvalid:       {"Signature is invalid or not found", "Подпись недействительна или не найдена"},
	msgPreview:                {"Preview", "Предпросмотр"},
	msgReady:                  {"Ready", "Готов"},
}

func tr(id messageID) string {
	pair, ok := messages[id]
	if !ok {
		return string(id)
	}
	if systemLanguageIsRussian() {
		return pair[1]
	}
	return pair[0]
}

func systemLanguageIsRussian() bool {
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANGUAGE", "LANG"} {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			continue
		}
		value = strings.ToLower(value)
		for _, part := range strings.FieldsFunc(value, func(r rune) bool {
			return r == ':' || r == ';' || r == ','
		}) {
			if strings.HasPrefix(strings.TrimSpace(part), "ru") {
				return true
			}
		}
	}
	return false
}
