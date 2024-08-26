package fsops

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"os"

	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

var (
	PostImgPath    = "./img/posts/"
	ProfilePicPath = "./img/profiles/"
)

type Certificate struct {
	CertPath string
	KeyPath  string
}

func MakePaths() *Certificate {
	cert := &Certificate{
		CertPath: "ssl/certificate.pem",
		KeyPath:  "ssl/key.pem",
	}
	return cert

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

func MakeImgDir() {
	if _, err := os.Stat("img"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("img", 0777)
		}
	}
	if _, err := os.Stat(PostImgPath); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(PostImgPath, 0777)
		}
	}
	if _, err := os.Stat(ProfilePicPath); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(ProfilePicPath, 0777)
		}
	}
}

func SaveImg(data []byte, path string, name string, extension string) {
	filePath := path + name + extension
	img, err := os.Create(filePath)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	defer img.Close()
	_, err = img.Write(data)
	if err != nil {
		logger.Error.Println(err)
		return
	}
}

func NameImg(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
