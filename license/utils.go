package license

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// ...

// DecodePrivateKey is a ...
func DecodePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKeyParsed, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try to parse the private key as PKCS#8 if PKCS#1 parsing fails
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error parsing private key: %v", err)
		}

		var ok bool
		privateKeyParsed, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not an RSA private key")
		}
	}

	return privateKeyParsed, nil
}

// DecodePublicKey decodes a public key from a byte slice
func DecodePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %v", err)
	}

	publicKeyParsed, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not an RSA public key")
	}

	return publicKeyParsed, nil
}

// KeyPair is a ...
type KeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
}

// KeyPairGenerate is a ...
func KeyPairGenerate() *KeyPair {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return &KeyPair{
		PublicKey:  publicKeyPEM,
		PrivateKey: privateKeyPEM,
	}
}
