package signature

import (
	"crypto/rsa"
	"errors"
)

// Signature to hold that signature needs, and contain parsed public and private key
type Signature struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	padding    PaddingDecider
}

// Options are needs to init the signature
type Options struct {
	PrivateKeyString string
	PublicKeyString  string
	PaddingType      PaddingType
}

// Init to init signature
func Init(opts Options) *Signature {
	privKey, err := parseRsaPrivateKeyFromPemStr(opts.PrivateKeyString)
	if err != nil {
		privKey = nil
	}

	publicKey, err := parsePublicKey(opts.PublicKeyString)
	if err != nil {
		publicKey = nil
	}

	return &Signature{
		privateKey: privKey,
		publicKey:  publicKey,
		padding:    decidePadding(opts.PaddingType),
	}
}

// Verify will return nil error if message and signature is match and verify
func (s *Signature) Verify(msg, signature string) error {
	if s == nil {
		return errors.New("signature is not init properly")
	}
	if s.publicKey == nil {
		return errors.New("public key is not setted or incorrect")
	}

	return s.padding.Verify(s.publicKey, msg, signature)
}

// Sign will return generated signature
func (s *Signature) Sign(msg []byte) (string, error) {
	if s == nil {
		return "", errors.New("signature is not init properly")
	}
	if s.privateKey == nil {
		return "", errors.New("private key is not setted or incorrect")
	}

	return s.padding.Sign(s.privateKey, msg)
}
