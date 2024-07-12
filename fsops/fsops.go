package fsops

import (
	"crypto/tls"
	"os"
)

type Certificate struct {
	CertPath string
	KeyPath  string
}

func MakePaths() (*Certificate, error) {
	HomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cert := &Certificate{
		CertPath: HomeDir + "/security/openssl/certificate.pem",
		KeyPath:  HomeDir + "/security/openssl/key.pem",
	}
	return cert, nil

}

func LoadCertificates(cert *Certificate) (*tls.Config, error) {
	certificates, err := tls.LoadX509KeyPair(cert.CertPath, cert.KeyPath)
	if err != nil {
		return nil, err
	}
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{certificates},
	}
	return tlsConf, nil
}
