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
	msgReason:                 {"Reason", "Причина"},
	msgDefaultReason:          {"Document signing", "Подписание документа"},
	msgSaveNextToSource:       {"Save next to source PDF", "Сохранять рядом с исходным PDF"},
	msgDetachedModeNote:       {"Linux mode creates a detached .sig for the stamped PDF.", "Linux-режим создает открепленную .sig для PDF со штампом."},
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
	msgBatchOutputNote:        {"Outputs use the _stamped suffix. The .sig file is created for each stamped PDF.", "Выходные файлы получают суффикс _stamped. Для каждого PDF со штампом создается .sig."},
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
