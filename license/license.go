package license

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	// ErrInvalidSignature is a ...
	ErrInvalidSignature = errors.New("Invalid signature")

	// ErrMalformedLicense is a ...
	ErrMalformedLicense = errors.New("Malformed License")

	// Read the private and public keys from the PEM files
	privateKey = readKeyFromFile("private_key.pem")
	publicKey  = readKeyFromFile("public_key.pem")
)

// Read the key from a PEM file
func readKeyFromFile(filename string) []byte {
	key, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading key from file %s: %v\n", filename, err)
		os.Exit(1)
	}
	return key
}

// License is a ...
type License struct {
	Iss string          `json:"iss,omitempty"` // Issued By
	Cus string          `json:"cus,omitempty"` // Customer ID
	Sub uint32          `json:"sub,omitempty"` // Subscriber ID
	Typ string          `json:"typ,omitempty"` // License Type
	Lim Limits          `json:"lim,omitempty"` // License Limit (e.g. Site)
	Iat int64           `json:"iat,omitempty"` // Issued At (timestamp)
	Exp int64           `json:"exp,omitempty"` // Expires At (timestamp)
	Dat json.RawMessage `json:"dat,omitempty"` // Metadata
}

// Limits is a ...
type Limits struct {
	Servers   int `json:"servers"`
	Companies int `json:"companies"`
	Users     int `json:"users"`
}

// Expired is a ...
func (l *License) Expired() bool {
	return l.Exp != 0 && time.Now().Unix() > l.Exp
}

// Encode is a ...
func (l *License) Encode(privateKey *rsa.PrivateKey) ([]byte, error) {
	msg, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	hash.Write(msg)
	digest := hash.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.Write(sig)
	buf.Write(msg)

	block := &pem.Block{
		Type:  "LICENSE KEY",
		Bytes: buf.Bytes(),
	}
	return pem.EncodeToMemory(block), nil
}

// Decode is a ...
func Decode(data []byte, publicKey *rsa.PublicKey) (*License, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrMalformedLicense
	}

	// Calculate the signature size based on the length of the block.Bytes
	sigSize := len(block.Bytes) - len(data) + len(block.Headers["Proc-Type"]) + len(block.Headers["DEK-Info"]) + 2
	sig := block.Bytes[:sigSize]
	msg := block.Bytes[sigSize:]

	hash := sha256.New()
	hash.Write(msg)
	digest := hash.Sum(nil)

	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest, sig)
	if err != nil {
		return nil, ErrInvalidSignature
	}

	out := new(License)
	err = json.Unmarshal(msg, out)
	return out, err
}

// GetPrivateKey is a ...
func GetPrivateKey() *rsa.PrivateKey {
	key, _ := DecodePrivateKey(privateKey)
	return key
}

// GetPublicKey is a ...
func GetPublicKey() *rsa.PublicKey {
	key, _ := DecodePublicKey(publicKey)
	return key
}
