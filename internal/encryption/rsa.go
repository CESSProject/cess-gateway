package encryption

import (
	"cess-httpservice/configs"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func Init() {
	if err := generateRSAKeyfile(2048); err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
		os.Exit(1)
	}
}

// generate key file
func generateRSAKeyfile(bits int) error {
	var err error
	_, err1 := os.Stat(configs.PrivateKeyfile)
	_, err2 := os.Stat(configs.PublicKeyfile)
	if err1 == nil && err2 == nil {
		return nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateBlock := pem.Block{Type: "RSA PRIVATE KEY", Bytes: X509PrivateKey}

	privateFile, err := os.Create(configs.PrivateKeyfile)
	if err != nil {
		return err
	}
	defer privateFile.Close()

	err = pem.Encode(privateFile, &privateBlock)
	if err != nil {
		return err
	}

	publicKey := privateKey.PublicKey
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	publicBlock := pem.Block{Type: "RSA PUBLIC KEY", Bytes: X509PublicKey}

	publicFile, err := os.Create(configs.PublicKeyfile)
	if err != nil {
		return err
	}
	defer publicFile.Close()
	err = pem.Encode(publicFile, &publicBlock)
	if err != nil {
		return err
	}
	return nil
}

// Parse private key file
func GetRSAPrivateKey(path string) *rsa.PrivateKey {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	block, _ := pem.Decode(buf)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privateKey
}

// Parse public key file
func GetRSAPublicKey(path string) *rsa.PublicKey {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	block, _ := pem.Decode(buf)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	return publicKey
}

// Parse private key
func ParsePrivateKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privateKey, err
}

// Parse public key
func ParsePublicKey(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	return publicKey, nil
}

// Calculate the signature
func CalcSign(msg []byte, privkey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.New()
	hash.Write(msg)
	bytes := hash.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, privkey, crypto.SHA256, bytes)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

// Verify signature
func VerifySign(msg []byte, sign []byte, pubkey *rsa.PublicKey) bool {
	hash := sha256.New()
	hash.Write(msg)
	bytes := hash.Sum(nil)
	err := rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, bytes, sign)
	return err == nil
}
