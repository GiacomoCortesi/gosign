package service

import (
	"encoding/base64"
	"fmt"

	"github.com/GiacomoCortesi/gosign/crypto"
	"github.com/GiacomoCortesi/gosign/domain"
	"github.com/google/uuid"
)

type signatureDeviceService struct {
	signatureDeviceRepository domain.SignatureDeviceRepository
	signerFactory             crypto.SignerFactory
}

// NewSignatureDeviceService return a SignatureDeviceService implementation
func NewSignatureDeviceService(repository domain.SignatureDeviceRepository) domain.SignatureDeviceService {
	return signatureDeviceService{
		signatureDeviceRepository: repository,
		signerFactory:             crypto.NewSignerFactory(),
	}
}

func generateKeyPair(a crypto.SignatureAlgorithm) (public, private []byte, err error) {
	switch a {
	case crypto.SignatureAlgorithmRSA:
		generator := crypto.RSAGenerator{}
		kp, err := generator.Generate()
		if err != nil {
			return public, private, err
		}
		public, private, err = crypto.NewRSAMarshaler().Marshal(*kp)
		if err != nil {
			return private, public, err
		}
	case crypto.SignatureAlgorithmECC:
		generator := crypto.ECCGenerator{}
		kp, err := generator.Generate()
		if err != nil {
			return public, private, err
		}
		public, private, err = crypto.NewECCMarshaler().Encode(*kp)
		if err != nil {
			return private, public, err
		}
	}
	return
}

// Create creates and return a new signature device
// If no ID is specified in the request, the ID is randomly generated
func (s signatureDeviceService) Create(sdreq domain.SignatureDeviceRequest) (domain.SignatureDeviceResponse, error) {
	// create random ID if not provided in request
	if sdreq.ID == "" {
		sdreq.ID = uuid.NewString()
	}
	public, private, err := generateKeyPair(sdreq.Algorithm)
	if err != nil {
		return domain.SignatureDeviceResponse{}, err
	}
	sdreq.PublicKey = public
	sdreq.PrivateKey = private
	return s.signatureDeviceRepository.Create(sdreq)
}

// Get retrieves a signature device given its ID
func (s signatureDeviceService) Get(deviceId string) (domain.SignatureDeviceResponse, error) {
	return s.signatureDeviceRepository.Get(deviceId)
}

// GetAll retrieves all available signature devices
func (s signatureDeviceService) GetAll() ([]domain.SignatureDeviceResponse, error) {
	return s.signatureDeviceRepository.GetAll()
}

// SignTransaction return the signed transaction data as a domain.SignatureResponse.
// Input data is extended to have this format: <signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded | device_id_base64_encoded>
// and then signed with appropriate algorithm
// After the signature has been created, the signature's counter value is incremented.
func (s signatureDeviceService) SignTransaction(deviceId string, data string) (domain.SignatureResponse, error) {
	// fetch the signature device from repository
	sdr, err := s.signatureDeviceRepository.Get(deviceId)
	if err != nil {
		return domain.SignatureResponse{}, err
	}

	// instantiate the appropriate signer for the device
	signer, err := s.signerFactory.CreateSigner(sdr.Algorithm, sdr.PrivateKey)
	if err != nil {
		return domain.SignatureResponse{}, err
	}

	// extend raw data:
	// <signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>
	var lastSignature string
	if sdr.SignatureCounter.Value() == 0 {
		lastSignature = sdr.ID
	} else {
		// Set lastSignature to the latest signature
		signatures, err := s.signatureDeviceRepository.GetAllSignature(deviceId)
		if err != nil {
			return domain.SignatureResponse{}, err
		}
		lastSignature = signatures[len(signatures)-1].Signature
	}
	lastSignature = base64.StdEncoding.EncodeToString([]byte(lastSignature))

	securedDataToBeSigned := fmt.Sprintf("%d_%s_%s", sdr.SignatureCounter, data, lastSignature)

	// sign the data
	signedData, err := signer.Sign([]byte(securedDataToBeSigned))
	if err != nil {
		return domain.SignatureResponse{}, err
	}

	sres := domain.SignatureResponse{
		Signature:  base64.StdEncoding.EncodeToString(signedData),
		SignedData: securedDataToBeSigned,
	}
	// add signature data to signature device
	if _, err = s.signatureDeviceRepository.AddSignature(deviceId, sres); err != nil {
		return domain.SignatureResponse{}, err
	}

	return sres, nil
}

// GetAllSignature return a slice of domain.SignatureResponse
func (s signatureDeviceService) GetAllSignature(deviceId string) ([]domain.SignatureResponse, error) {
	return s.signatureDeviceRepository.GetAllSignature(deviceId)
}
