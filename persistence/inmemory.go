/*
Package persistence implements domain repository interface for the gosign microservice

In memory persistence implementation use mutex to allow concurrent client access to the data
*/
package persistence

import (
	"sync"

	"github.com/GiacomoCortesi/gosign/domain"
)

type inMemorySignatureDeviceRepository struct {
	signatureDevice  map[string]domain.SignatureDeviceResponse
	deviceSignatures map[string][]domain.SignatureResponse

	mu sync.RWMutex
}

// NewInMemorySignatureDeviceRepository return an in memory implementation of the
// domain.SignatureDeviceRepository interface
func NewInMemorySignatureDeviceRepository() domain.SignatureDeviceRepository {
	return &inMemorySignatureDeviceRepository{
		signatureDevice:  make(map[string]domain.SignatureDeviceResponse),
		deviceSignatures: make(map[string][]domain.SignatureResponse),
		mu:               sync.RWMutex{},
	}
}

// Create create a new signature device
func (r *inMemorySignatureDeviceRepository) Create(sdreq domain.SignatureDeviceRequest) (sdres domain.SignatureDeviceResponse, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.signatureDevice[sdreq.ID]; exist {
		return sdres, domain.ErrSignatureDeviceAlreadyExist
	}
	sdres = domain.SignatureDeviceResponse{
		ID:               sdreq.ID,
		Algorithm:        sdreq.Algorithm,
		Label:            sdreq.Label,
		SignatureCounter: 0,
		PrivateKey:       sdreq.PrivateKey,
		PublicKey:        sdreq.PublicKey,
	}
	r.signatureDevice[sdreq.ID] = sdres
	r.deviceSignatures[sdreq.ID] = make([]domain.SignatureResponse, 0)
	return
}

// AddSignature add a new signature to the signature device and updates the signature counter
func (r *inMemorySignatureDeviceRepository) AddSignature(deviceId string, sres domain.SignatureResponse) (sdres domain.SignatureDeviceResponse, err error) {
	sdres, err = r.Get(deviceId)
	if err != nil {
		return sdres, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	sdres.SignatureCounter.Increment()
	r.signatureDevice[deviceId] = sdres
	r.deviceSignatures[deviceId] = append(r.deviceSignatures[deviceId], sres)
	return
}

// GetAllSignature return all available signatures for the specified device
func (r *inMemorySignatureDeviceRepository) GetAllSignature(deviceId string) (sres []domain.SignatureResponse, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sres, exist := r.deviceSignatures[deviceId]
	if !exist {
		return sres, domain.ErrSignatureDeviceNotFound
	}
	return
}

// GetAll return all available signature devices
func (r *inMemorySignatureDeviceRepository) GetAll() ([]domain.SignatureDeviceResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sdresList := []domain.SignatureDeviceResponse{}
	for _, sdres := range r.signatureDevice {
		sdresList = append(sdresList, sdres)
	}
	return sdresList, nil
}

// Get return the signature device having the specified ID
func (r *inMemorySignatureDeviceRepository) Get(deviceId string) (sdres domain.SignatureDeviceResponse, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sdres, exist := r.signatureDevice[deviceId]
	if !exist {
		return sdres, domain.ErrSignatureDeviceNotFound
	}
	return
}
