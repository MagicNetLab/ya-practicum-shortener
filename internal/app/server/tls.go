package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"time"
)

func hasTLSCertsExists(certFileName string, keyFileName string) bool {
	if _, err := os.Stat(certFileName); errors.Is(err, os.ErrNotExist) {
		return false
	}
	if _, err := os.Stat(keyFileName); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// создание сертификата для https сервера
func createTLSCerts(certFileName string, keyFileName string) error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Roga&Kopita"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	err = writeCertToFile(certFileName, "CERTIFICATE", certBytes)
	if err != nil {
		return err
	}

	err = writeCertToFile(keyFileName, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))
	if err != nil {
		return err
	}

	return nil
}

// сохраняет сертификат на диск
func writeCertToFile(filepath string, cypherType string, cypher []byte) error {
	var (
		buf bytes.Buffer
		f   *os.File
	)
	err := pem.Encode(&buf, &pem.Block{
		Type:  cypherType,
		Bytes: cypher,
	})
	if err != nil {
		return err
	}

	f, err = os.Create(filepath)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(f)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
