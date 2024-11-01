package signature

import (
	"crypto/rsa"
)

// PaddingType is supported padding
type PaddingType int

const (
	// PaddingTypePSS PSS padding, the most secure and the default
	PaddingTypePSS PaddingType = 0

	// PaddingTypePKCS1v15 is a PKCS1v15 constant
	PaddingTypePKCS1v15 PaddingType = 1
)

type paddingPSS struct{}
type paddingPKCS1v15 struct{}

// PaddingDecider padding decider
type PaddingDecider interface {
	// Verify verifies the `signature` using selected padding
	Verify(pubKey *rsa.PublicKey, msg, signature string) error
	// Sign signs the `msg` using selected padding
	Sign(privKey *rsa.PrivateKey, msg []byte) (string, error)
}

func decidePadding(padtype PaddingType) PaddingDecider {
	switch padtype {
	case PaddingTypePSS:
		return &paddingPSS{}
	case PaddingTypePKCS1v15:
		return &paddingPKCS1v15{}
	default:
		return &paddingPSS{}
	}
}
