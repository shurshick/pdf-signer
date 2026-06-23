package main

import (
	"fmt"
	"strings"
)

func FriendlyErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	msg := err.Error()
	lower := strings.ToLower(msg)

	switch {
	case strings.Contains(lower, "certmgr") && strings.Contains(lower, "not found"):
		return tr(msgCertmgrNotFound)
	case strings.Contains(lower, "csptest") && strings.Contains(lower, "not found"):
		return tr(msgCsptestNotFound)
	case strings.Contains(lower, "no certificates") || strings.Contains(lower, "сертификаты") && strings.Contains(lower, "не найдены"):
		return tr(msgNoCertsFound)
	case strings.Contains(lower, "certificate cn is empty") || strings.Contains(lower, "пустой cn"):
		return tr(msgCertCNEmpty)
	case strings.Contains(lower, "pdf is not specified") || strings.Contains(lower, "не указан pdf"):
		return tr(msgPDFNotSpecified)
	case strings.Contains(lower, "signing error") || strings.Contains(lower, "ошибка подписи"):
		return tr(msgSigningError)
	case strings.Contains(lower, "embedded signing error") || strings.Contains(lower, "ошибка встроенного подписания"):
		return tr(msgEmbeddedSignError)
	case strings.Contains(lower, "stamp") && strings.Contains(lower, "error"):
		return tr(msgStampError)
	case strings.Contains(lower, "permission denied") || strings.Contains(lower, "access denied"):
		return tr(msgPermissionDenied)
	case strings.Contains(lower, "expired") || strings.Contains(lower, "истек"):
		return tr(msgCertExpired)
	case strings.Contains(lower, "not valid yet") || strings.Contains(lower, "еще не действителен"):
		return tr(msgCertNotValidYet)
	case strings.Contains(lower, "private key") || strings.Contains(lower, "закрытый ключ"):
		return tr(msgPrivateKeyMissing)
	case strings.Contains(lower, "crypto") && strings.Contains(lower, "error"):
		return tr(msgCryptoError)
	default:
		return fmt.Sprintf("%s: %s", tr(msgError), msg)
	}
}
