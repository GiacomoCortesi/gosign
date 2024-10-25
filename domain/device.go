/*
Package domain provides the business logic and interfaces required by the gosign microservice:
 - signature device service definition and implementation
 - signature device repository interface definition for data access layer operations
 - common types and errors
*/

package domain

import (
	"errors"
	"sync/atomic"

	"github.com/GiacomoCortesi/gosign/crypto"
)

// Signature device custom errors
var (
	ErrSignatureDeviceNotFound     = errors.New("signature device not found")
	ErrSignatureDeviceAlreadyExist = errors.New("signature device already exist")
)

// SignatureDeviceRepository provides methods for performing data access layer operations
// on signature devices
type SignatureDeviceRepository interface {
	Create(SignatureDeviceRequest) (SignatureDeviceResponse, error)
	GetAll() ([]SignatureDeviceResponse, error)
	Get(deviceId string) (SignatureDeviceResponse, error)
	AddSignature(deviceId string, sres SignatureResponse) (SignatureDeviceResponse, error)
	GetAllSignature(deviceId string) ([]SignatureResponse, error)
}

// SignatureDeviceService provide methods for managing signature devices
type SignatureDeviceService interface {
	Create(SignatureDeviceRequest) (SignatureDeviceResponse, error)
	GetAll() ([]SignatureDeviceResponse, error)
	Get(deviceId string) (SignatureDeviceResponse, error)
	SignTransaction(deviceId string, data string) (SignatureResponse, error)
	GetAllSignature(deviceId string) ([]SignatureResponse, error)
}

// SignatureDeviceRequest represent a signature device request
type SignatureDeviceRequest struct {
	ID         string                    `json:"id"`
	Algorithm  crypto.SignatureAlgorithm `json:"algorithm"`
	Label      string                    `json:"label,omitempty"`
	PrivateKey []byte                    `json:"-"`
	PublicKey  []byte                    `json:"-"`
}

// SignatureDeviceResponse represent a signature device response
type SignatureDeviceResponse struct {
	ID               string                    `json:"id"`
	Algorithm        crypto.SignatureAlgorithm `json:"algorithm"`
	Label            string                    `json:"label,omitempty"`
	SignatureCounter SignatureCounter          `json:"signature_counter"`
	PrivateKey       []byte                    `json:"-"`
	PublicKey        []byte                    `json:"-"`
}

// SignatureRequest represent the device sign transaction request
type SignatureRequest struct {
	Data string `json:"data"`
}

// SignatureResponse represent the device sign transaction response
type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

// SignatureCounter represent a thread-safe integer counter
type SignatureCounter int64

// Increment increments the counter
func (c *SignatureCounter) Increment() {
	atomic.AddInt64((*int64)(c), 1)
}

// Value return the int64 value of the counter
func (c *SignatureCounter) Value() int64 {
	return atomic.LoadInt64((*int64)(c))
}
