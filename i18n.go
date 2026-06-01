package main

import (
	"os"
	"strings"
)

type messageID string

const (
	msgWindowTitle       messageID = "window_title"
	msgPDFNotSelected    messageID = "pdf_not_selected"
	msgOutPDFNotSelected messageID = "out_pdf_not_selected"
	msgCertNotSelected   messageID = "cert_not_selected"
	msgOwner             messageID = "owner"
	msgIssuer            messageID = "issuer"
	msgSelectPDF         messageID = "select_pdf"
	msgSavePDFAs         messageID = "save_pdf_as"
	msgSignAndStamp      messageID = "sign_and_stamp"
	msgError             messageID = "error"
	msgChoosePDF         messageID = "choose_pdf"
	msgChooseCert        messageID = "choose_cert"
	msgChooseOutPDF      messageID = "choose_out_pdf"
	msgDone              messageID = "done"
	msgSignature         messageID = "signature"
	msgStampedPDF        messageID = "stamped_pdf"
	msgCertificate       messageID = "certificate"
	msgScale             messageID = "scale"
	msgCertmgrError      messageID = "certmgr_error"
	msgNoCerts           messageID = "no_certs"
	msgNoPDF             messageID = "no_pdf"
	msgAbsPDFError       messageID = "abs_pdf_error"
	msgEmptyCertCN       messageID = "empty_cert_cn"
	msgSignError         messageID = "sign_error"
	msgSignatureMissing  messageID = "signature_missing"
	msgInputPDFMissing   messageID = "input_pdf_missing"
	msgOutputPDFMissing  messageID = "output_pdf_missing"
	msgStampFileMissing  messageID = "stamp_file_missing"
	msgStampPDFError     messageID = "stamp_pdf_error"
	msgStampTitle        messageID = "stamp_title"
	msgDate              messageID = "date"
	msgSerialNumber      messageID = "serial_number"
	msgSignatureShort    messageID = "signature_short"
)

var messages = map[messageID][2]string{
	msgWindowTitle:       {"PDF signer + stamp", "PDF подпись + штамп"},
	msgPDFNotSelected:    {"PDF is not selected", "PDF не выбран"},
	msgOutPDFNotSelected: {"Output PDF is not selected", "Выходной PDF не выбран"},
	msgCertNotSelected:   {"Certificate is not selected", "Сертификат не выбран"},
	msgOwner:             {"Owner", "Владелец"},
	msgIssuer:            {"Issuer", "Издатель"},
	msgSelectPDF:         {"Select PDF", "Выбрать PDF"},
	msgSavePDFAs:         {"Save PDF as", "Сохранить PDF как"},
	msgSignAndStamp:      {"Sign and stamp", "Подписать и поставить штамп"},
	msgError:             {"Error", "Ошибка"},
	msgChoosePDF:         {"Select a PDF file", "Выберите PDF"},
	msgChooseCert:        {"Select a certificate", "Выберите сертификат"},
	msgChooseOutPDF:      {"Choose output PDF", "Укажите выходной PDF"},
	msgDone:              {"Done", "Готово"},
	msgSignature:         {"Signature", "Подпись"},
	msgStampedPDF:        {"Stamped PDF", "PDF со штампом"},
	msgCertificate:       {"Certificate", "Сертификат"},
	msgScale:             {"Scale", "Масштаб"},
	msgCertmgrError:      {"certmgr error", "ошибка certmgr"},
	msgNoCerts:           {"no certificates found in uMy store", "сертификаты в хранилище uMy не найдены"},
	msgNoPDF:             {"PDF is not specified", "не указан PDF"},
	msgAbsPDFError:       {"failed to resolve absolute PDF path", "не удалось получить абсолютный путь к PDF"},
	msgEmptyCertCN:       {"certificate CN is empty", "у сертификата пустой CN"},
	msgSignError:         {"signing error", "ошибка подписи"},
	msgSignatureMissing:  {"signature file was not created", "файл подписи не создан"},
	msgInputPDFMissing:   {"input PDF is not specified", "не указан входной PDF"},
	msgOutputPDFMissing:  {"output PDF is not specified", "не указан выходной PDF"},
	msgStampFileMissing:  {"stamp file is not specified", "не указан файл штампа"},
	msgStampPDFError:     {"failed to add stamp to PDF", "ошибка добавления штампа в PDF"},
	msgStampTitle:        {"ELECTRONIC SIGNATURE", "ЭЛЕКТРОННАЯ ПОДПИСЬ"},
	msgDate:              {"Date", "Дата"},
	msgSerialNumber:      {"Serial number", "Серийный номер"},
	msgSignatureShort:    {"ES", "ЭП"},
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
