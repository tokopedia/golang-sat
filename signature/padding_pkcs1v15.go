package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

// Verify will verify the message using PKCS1v15 padding
func (p *paddingPKCS1v15) Verify(pubKey *rsa.PublicKey, msg, signature string) error {
	message := []byte(msg)
	bSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256(message)
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], bSignature)
}

// Sign will generate signature to the message using PKCS1v15 padding
func (p *paddingPKCS1v15) Sign(privKey *rsa.PrivateKey, msg []byte) (string, error) {
	hashed := sha256.Sum256(msg)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
