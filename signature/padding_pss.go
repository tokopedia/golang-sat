package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

// Verify will verify the message using PSS padding
func (p *paddingPSS) Verify(pubKey *rsa.PublicKey, msg, signature string) error {
	if strings.TrimSpace(signature) == "" {
		return errors.New("signature is empty")
	}

	message := []byte(msg)
	bSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256(message)
	return rsa.VerifyPSS(pubKey, crypto.SHA256, hashed[:], bSignature, nil)
}

// Sign will generate signature to the message using PSS padding
func (p *paddingPSS) Sign(privKey *rsa.PrivateKey, msg []byte) (string, error) {
	hashed := sha256.Sum256(msg)
	signature, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, hashed[:], nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
