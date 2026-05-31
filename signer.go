package main

type SignResult struct {
    SignaturePath string
}

type Signer interface {
    SignFile(pdfPath string, cert CertInfo) (SignResult, error)
}