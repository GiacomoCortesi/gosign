package service

import (
	"reflect"
	"testing"

	"github.com/GiacomoCortesi/gosign/crypto"
	"github.com/GiacomoCortesi/gosign/domain"
	"github.com/GiacomoCortesi/gosign/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_signatureDeviceService_SignTransaction_DeviceID(t *testing.T) {
	mockSignerFactory := mocks.MockSignerFactory{}
	mockSigner := mocks.MockSigner{}
	mockSigner.On("Sign", mock.Anything).Return([]byte("thesignature"), nil)
	mockSignerFactory.On("CreateSigner", mock.Anything).Return(mockSigner, nil)

	mockRepository := mocks.MockSignatureDeviceRepository{}
	mockRepository.On("Get", mock.Anything).Return(domain.SignatureDeviceResponse{
		ID:               "someid",
		Algorithm:        crypto.SignatureAlgorithmRSA,
		SignatureCounter: 0,
	}, nil)
	mockRepository.On("AddSignature", mock.Anything, mock.Anything).Return(domain.SignatureDeviceResponse{
		ID:               "someid",
		Algorithm:        crypto.SignatureAlgorithmRSA,
		SignatureCounter: 1,
	}, nil)
	type fields struct {
		signatureDeviceRepository domain.SignatureDeviceRepository
	}
	type args struct {
		deviceId string
		data     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.SignatureResponse
		wantErr bool
	}{
		{
			name: "sign transaction success - last signature with encoded device ID",
			fields: fields{
				signatureDeviceRepository: mockRepository,
			},
			args: args{
				deviceId: "someid",
				data:     "somedata",
			},
			want: domain.SignatureResponse{
				Signature:  "dGhlc2lnbmF0dXJl",
				SignedData: "0_somedata_c29tZWlk",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := signatureDeviceService{
				signatureDeviceRepository: tt.fields.signatureDeviceRepository,
				signerFactory:             &mockSignerFactory,
			}
			got, err := s.SignTransaction(tt.args.deviceId, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("signatureDeviceService.SignTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("signatureDeviceService.SignTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_signatureDeviceService_SignTransaction_LastSignature(t *testing.T) {
	mockSignerFactory := mocks.MockSignerFactory{}
	mockSigner := mocks.MockSigner{}
	mockSigner.On("Sign", mock.Anything).Return([]byte("thesignature"), nil)
	mockSignerFactory.On("CreateSigner", mock.Anything).Return(mockSigner, nil)

	mockRepository := mocks.MockSignatureDeviceRepository{}
	mockRepository.On("Get", mock.Anything).Return(domain.SignatureDeviceResponse{
		ID:               "someid",
		Algorithm:        crypto.SignatureAlgorithmRSA,
		SignatureCounter: 1,
	}, nil)
	mockRepository.On("GetAllSignature", mock.Anything).Return([]domain.SignatureResponse{
		{
			Signature:  "cHJldmlvdXNzaWduYXR1cmUK",
			SignedData: "0_somepreviouslysigneddata_c29tZWlk",
		},
	}, nil)
	mockRepository.On("AddSignature", mock.Anything, mock.Anything).Return(domain.SignatureDeviceResponse{
		ID:               "someid",
		Algorithm:        crypto.SignatureAlgorithmRSA,
		SignatureCounter: 2,
	}, nil)
	type fields struct {
		signatureDeviceRepository domain.SignatureDeviceRepository
	}
	type args struct {
		deviceId string
		data     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.SignatureResponse
		wantErr bool
	}{
		{
			name: "sign transaction success - last signature",
			fields: fields{
				signatureDeviceRepository: mockRepository,
			},
			args: args{
				deviceId: "someid",
				data:     "somedata",
			},
			want: domain.SignatureResponse{
				Signature:  "dGhlc2lnbmF0dXJl",
				SignedData: "1_somedata_Y0hKbGRtbHZkWE56YVdkdVlYUjFjbVVL",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := signatureDeviceService{
				signatureDeviceRepository: tt.fields.signatureDeviceRepository,
				signerFactory:             &mockSignerFactory,
			}
			got, err := s.SignTransaction(tt.args.deviceId, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("signatureDeviceService.SignTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("signatureDeviceService.SignTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
