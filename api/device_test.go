package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GiacomoCortesi/gosign/crypto"
	"github.com/GiacomoCortesi/gosign/domain"
	"github.com/GiacomoCortesi/gosign/mocks"
	"github.com/stretchr/testify/mock"
)

func TestServer_GetAllSignatureDevice(t *testing.T) {
	mockServiceNoDevices := mocks.MockSignatureDeviceService{}
	mockServiceNoDevices.On("GetAll").Return([]domain.SignatureDeviceResponse{}, nil)

	mockServiceWithDevices := mocks.MockSignatureDeviceService{}
	mockServiceWithDevices.On("GetAll").Return([]domain.SignatureDeviceResponse{
		{
			ID:               "someid",
			Label:            "somelabel",
			Algorithm:        crypto.SignatureAlgorithmRSA,
			SignatureCounter: 5,
		},
	}, nil)

	type fields struct {
		listenAddress          string
		signatureDeviceService domain.SignatureDeviceService
	}
	tests := []struct {
		name       string
		fields     fields
		want       Response
		wantStatus int
	}{
		{
			name: "get all devices handler success - no devices",
			fields: fields{
				listenAddress:          "8080",
				signatureDeviceService: mockServiceNoDevices,
			},
			want:       Response{Data: []domain.SignatureDeviceResponse{}},
			wantStatus: http.StatusOK,
		},
		{
			name: "get all devices handler success - with devices",
			fields: fields{
				listenAddress:          "8080",
				signatureDeviceService: mockServiceWithDevices,
			},
			want: Response{
				Data: []domain.SignatureDeviceResponse{
					{
						ID:               "someid",
						Label:            "somelabel",
						Algorithm:        crypto.SignatureAlgorithmRSA,
						SignatureCounter: 5,
					},
				},
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				listenAddress:          tt.fields.listenAddress,
				signatureDeviceService: tt.fields.signatureDeviceService,
			}
			testServer := httptest.NewServer(http.HandlerFunc(s.GetAllSignatureDevice))
			defer testServer.Close()
			resp, err := http.Get(testServer.URL)
			if err != nil {
				t.Error(err)
			}
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("want status %d but got %d", tt.wantStatus, resp.StatusCode)
			}
			mockService := tt.fields.signatureDeviceService.(mocks.MockSignatureDeviceService)
			mockService.AssertExpectations(t)
		})
	}
}

func TestServer_GetSignatureDevice(t *testing.T) {
	mockServiceNoDevice := mocks.MockSignatureDeviceService{}
	mockServiceNoDevice.On("Get", mock.Anything).Return(domain.SignatureDeviceResponse{}, domain.ErrSignatureDeviceNotFound)

	mockServiceWithDevice := mocks.MockSignatureDeviceService{}
	mockServiceWithDevice.On("Get", mock.Anything).Return(domain.SignatureDeviceResponse{
		ID:               "someid",
		Label:            "somelabel",
		Algorithm:        crypto.SignatureAlgorithmRSA,
		SignatureCounter: 5,
	}, nil)

	type fields struct {
		listenAddress          string
		signatureDeviceService domain.SignatureDeviceService
	}
	tests := []struct {
		name       string
		fields     fields
		want       Response
		wantStatus int
	}{
		{
			name: "get device handler success",
			fields: fields{
				listenAddress:          "8080",
				signatureDeviceService: mockServiceWithDevice,
			},
			want: Response{
				Data: domain.SignatureDeviceResponse{
					ID:               "someid",
					Label:            "somelabel",
					Algorithm:        crypto.SignatureAlgorithmRSA,
					SignatureCounter: 5,
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "get device handler failure - device missing",
			fields: fields{
				listenAddress:          "8080",
				signatureDeviceService: mockServiceNoDevice,
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				listenAddress:          tt.fields.listenAddress,
				signatureDeviceService: tt.fields.signatureDeviceService,
			}
			testServer := httptest.NewServer(http.HandlerFunc(s.GetSignatureDevice))
			defer testServer.Close()
			resp, err := http.Get(testServer.URL)
			if err != nil {
				t.Error(err)
			}
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("want status %d but got %d", tt.wantStatus, resp.StatusCode)
			}
			mockService := tt.fields.signatureDeviceService.(mocks.MockSignatureDeviceService)
			mockService.AssertExpectations(t)
		})
	}
}
