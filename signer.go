package main

type SignResult struct {
    SignaturePath string
    SignedPDFPath string
}

type SignMode int

const (
    SignModeEmbedded  SignMode = iota
    SignModeDetached
    SignModeBoth
)

type Signer interface {
    SignFile(pdfPath string, cert CertInfo) (SignResult, error)
    SignFileTo(pdfPath string, cert CertInfo, sigPath string) (SignResult, error)
    SignFileEmbedded(pdfPath string, cert CertInfo) (SignResult, error)
}