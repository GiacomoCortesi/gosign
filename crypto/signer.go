package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
)

var ErrInvalidSignatureAlgorithm = errors.New("invalid signature algorithm")

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// SignatureAlgorithm is an utility enum type identifying supported algorithms for signing data
type SignatureAlgorithm int

const (
	SignatureAlgorithmRSA SignatureAlgorithm = iota // Signature algorithm RSA
	SignatureAlgorithmECC                           // Signature algorithm ECC
)

// MarshalJSON encodes the SignatureAlgorithm as a string.
func (sa SignatureAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(sa.String())
}

// UnmarshalJSON decodes the SignatureAlgorithm from a string.
func (sa *SignatureAlgorithm) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "RSA":
		*sa = SignatureAlgorithmRSA
	case "ECC":
		*sa = SignatureAlgorithmECC
	default:
		return ErrInvalidSignatureAlgorithm
	}

	return nil
}

// String return the string representation of the signature algorithm
func (s SignatureAlgorithm) String() string {
	return []string{"RSA", "ECC"}[s]
}

// RSASigner implement Signer interface for RSA algorithm
type RSASigner struct {
	pk *rsa.PrivateKey
}

// Sign return the RSA signed data
func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashedDataToBeSigned := sha256.Sum256(dataToBeSigned)
	return rsa.SignPKCS1v15(rand.Reader, s.pk, crypto.SHA256, hashedDataToBeSigned[:])
}

// NewRSASigner return an RSASigner instance
func NewRSASigner(pk rsa.PrivateKey) (*RSASigner, error) {
	return &RSASigner{
		pk: &pk,
	}, nil
}

// ECCSigner implement Signer interface for ECC algorithm
type ECCSigner struct {
	pk *ecdsa.PrivateKey
}

// Sign return the ECC signed data
func (s *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashedDataToBeSigned := sha256.Sum256(dataToBeSigned)
	return ecdsa.SignASN1(rand.Reader, s.pk, hashedDataToBeSigned[:])
}

// NewECCSigner return an ECCSigner instance
func NewECCSigner(pk ecdsa.PrivateKey) (*ECCSigner, error) {
	return &ECCSigner{
		pk: &pk,
	}, nil
}

type SignerFactory interface {
	CreateSigner(a SignatureAlgorithm, privateKey []byte) (Signer, error)
}

type signerFactory struct{}

// CreateSigner is a factory for creating a Signer instance based on the specified signature algorithm
func (sf signerFactory) CreateSigner(a SignatureAlgorithm, pk []byte) (s Signer, err error) {
	switch a {
	case SignatureAlgorithmECC:
		kp, err := NewECCMarshaler().Decode(pk)
		if err != nil {
			return nil, err
		}
		return NewECCSigner(*kp.Private)
	case SignatureAlgorithmRSA:
		kp, err := NewRSAMarshaler().Unmarshal(pk)
		if err != nil {
			return nil, err
		}
		return NewRSASigner(*kp.Private)
	default:
		return nil, ErrInvalidSignatureAlgorithm
	}
}

func NewSignerFactory() SignerFactory {
	return signerFactory{}
}
