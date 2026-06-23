package main

import (
	"testing"
)

func TestTrReturnsNonEmpty(t *testing.T) {
	ids := []messageID{
		msgWindowTitle,
		msgPDFNotSelected,
		msgSignAndStamp,
		msgAbout,
		msgModeEmbedded,
		msgModeDetached,
		msgModeBoth,
		msgSigningMode,
		msgEmbeddedSignature,
	}
	for _, id := range ids {
		result := tr(id)
		if result == "" {
			t.Errorf("tr(%q) returned empty string", id)
		}
	}
}

func TestTrUnknownID(t *testing.T) {
	result := tr("nonexistent_key_12345")
	if result != "nonexistent_key_12345" {
		t.Errorf("tr(unknown) = %q, want %q", result, "nonexistent_key_12345")
	}
}

func TestModeLabels(t *testing.T) {
	embedded := tr(msgModeEmbedded)
	detached := tr(msgModeDetached)
	both := tr(msgModeBoth)

	if embedded == detached || embedded == both || detached == both {
		t.Error("mode labels should be unique")
	}
}

func TestSystemLanguageIsRussian(t *testing.T) {
	result := systemLanguageIsRussian()
	if result {
		t.Log("systemLanguageIsRussian() = true (Russian locale detected)")
	} else {
		t.Log("systemLanguageIsRussian() = false (non-Russian locale)")
	}
}

func TestAllMessageIDsHaveTranslations(t *testing.T) {
	allIDs := []messageID{
		msgWindowTitle, msgPDFNotSelected, msgOutPDFNotSelected,
		msgCertNotSelected, msgOwner, msgIssuer, msgSelectPDF,
		msgClearPDFs, msgSelectOutputFolder, msgBrowse, msgAbout,
		msgVersion, msgProjectLink, msgCheckUpdates,
		msgUpdateAvailableTitle, msgUpdateAvailableBody,
		msgNoUpdatesTitle, msgNoUpdatesBody, msgUpdateCheckFailed,
		msgCurrentVersion, msgLatestVersion, msgDownload, msgClose,
		msgReason, msgDefaultReason, msgSaveNextToSource,
		msgDetachedModeNote, msgSignAndStamp, msgError, msgChoosePDF,
		msgChooseCert, msgChooseOutPDF, msgDone, msgSignature,
		msgStampedPDF, msgSelectedPDFs, msgProcessedFiles,
		msgBatchOutputNote, msgCertificate, msgScale,
		msgCertmgrError, msgNoCerts, msgNoPDF, msgAbsPDFError,
		msgAbsSignatureError, msgEmptyCertCN, msgSignError,
		msgSignatureMissing, msgOutputSignatureMissing,
		msgInputPDFMissing, msgOutputPDFMissing, msgStampFileMissing,
		msgStampPDFError, msgStampTitle, msgDate, msgSerialNumber,
		msgSignatureShort, msgSigningMode, msgModeEmbedded,
		msgModeDetached, msgModeBoth, msgModeEmbeddedDesc,
		msgModeDetachedDesc, msgModeBothDesc, msgSignEmbeddedError,
		msgSignedPDF, msgEmbeddedSignature,
	}

	for _, id := range allIDs {
		if _, ok := messages[id]; !ok {
			t.Errorf("messageID %q not found in messages map", id)
		}
		en := messages[id][0]
		ru := messages[id][1]
		if en == "" {
			t.Errorf("English translation empty for %q", id)
		}
		if ru == "" {
			t.Errorf("Russian translation empty for %q", id)
		}
	}
}
