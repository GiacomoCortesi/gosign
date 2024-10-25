/*
Package mocks provide mocked interfaces to be used in unit tests
*/
package mocks

import (
	"github.com/GiacomoCortesi/gosign/crypto"
	"github.com/GiacomoCortesi/gosign/domain"
	"github.com/stretchr/testify/mock"
)

type MockSignerFactory struct {
	mock.Mock
}

func (m MockSignerFactory) CreateSigner(algo crypto.SignatureAlgorithm, pk []byte) (crypto.Signer, error) {
	args := m.Called(algo)
	return args.Get(0).(crypto.Signer), args.Error(1)
}

type MockSigner struct {
	mock.Mock
}

func (m MockSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	args := m.Called(dataToBeSigned)
	return args.Get(0).([]byte), args.Error(1)
}

type MockSignatureDeviceRepository struct {
	mock.Mock
}

func (m MockSignatureDeviceRepository) Create(req domain.SignatureDeviceRequest) (domain.SignatureDeviceResponse, error) {
	args := m.Called(req)
	return args.Get(0).(domain.SignatureDeviceResponse), args.Error(1)
}

func (m MockSignatureDeviceRepository) GetAll() ([]domain.SignatureDeviceResponse, error) {
	args := m.Called()
	return args.Get(0).([]domain.SignatureDeviceResponse), args.Error(1)
}

func (m MockSignatureDeviceRepository) Get(deviceId string) (domain.SignatureDeviceResponse, error) {
	args := m.Called(deviceId)
	return args.Get(0).(domain.SignatureDeviceResponse), args.Error(1)
}

func (m MockSignatureDeviceRepository) AddSignature(deviceId string, sres domain.SignatureResponse) (domain.SignatureDeviceResponse, error) {
	args := m.Called(deviceId, sres)
	return args.Get(0).(domain.SignatureDeviceResponse), args.Error(1)
}

func (m MockSignatureDeviceRepository) GetAllSignature(deviceId string) ([]domain.SignatureResponse, error) {
	args := m.Called(deviceId)
	return args.Get(0).([]domain.SignatureResponse), args.Error(1)
}

type MockSignatureDeviceService struct {
	mock.Mock
}

func (m MockSignatureDeviceService) Create(req domain.SignatureDeviceRequest) (domain.SignatureDeviceResponse, error) {
	args := m.Called(req)
	return args.Get(0).(domain.SignatureDeviceResponse), args.Error(1)
}

func (m MockSignatureDeviceService) GetAll() ([]domain.SignatureDeviceResponse, error) {
	args := m.Called()
	return args.Get(0).([]domain.SignatureDeviceResponse), args.Error(1)
}

func (m MockSignatureDeviceService) Get(deviceId string) (domain.SignatureDeviceResponse, error) {
	args := m.Called(deviceId)
	return args.Get(0).(domain.SignatureDeviceResponse), args.Error(1)
}

func (m MockSignatureDeviceService) SignTransaction(deviceId string, data string) (domain.SignatureResponse, error) {
	args := m.Called(deviceId, data)
	return args.Get(0).(domain.SignatureResponse), args.Error(1)
}

func (m MockSignatureDeviceService) GetAllSignature(deviceId string) ([]domain.SignatureResponse, error) {
	args := m.Called(deviceId)
	return args.Get(0).([]domain.SignatureResponse), args.Error(1)
}
