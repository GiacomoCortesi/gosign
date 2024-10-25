package persistence

import (
	"reflect"
	"sync"
	"testing"

	"github.com/GiacomoCortesi/gosign/crypto"
	"github.com/GiacomoCortesi/gosign/domain"
)

func TestNewInMemorySignatureDeviceRepository(t *testing.T) {
	tests := []struct {
		name string
		want domain.SignatureDeviceRepository
	}{
		{
			name: "valid in memory signature device repository creation",
			want: &inMemorySignatureDeviceRepository{
				signatureDevice:  make(map[string]domain.SignatureDeviceResponse),
				deviceSignatures: make(map[string][]domain.SignatureResponse),
				mu:               sync.RWMutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInMemorySignatureDeviceRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInMemorySignatureDeviceRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inMemorySignatureDeviceRepository_Create(t *testing.T) {
	type fields struct {
		signatureDevice  map[string]domain.SignatureDeviceResponse
		deviceSignatures map[string][]domain.SignatureResponse
	}
	type args struct {
		sdreq domain.SignatureDeviceRequest
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want domain.SignatureDeviceResponse
		wantErr   bool
	}{
		{
			name: "create signature device success",
			fields: fields{
				signatureDevice:  make(map[string]domain.SignatureDeviceResponse),
				deviceSignatures: make(map[string][]domain.SignatureResponse),
			},
			args: args{
				sdreq: domain.SignatureDeviceRequest{
					ID:        "someid",
					Label:     "some label",
					Algorithm: crypto.SignatureAlgorithmRSA,
				},
			},
			want: domain.SignatureDeviceResponse{
				ID:               "someid",
				Label:            "some label",
				Algorithm:        crypto.SignatureAlgorithmRSA,
				SignatureCounter: 0,
			},
			wantErr: false,
		},
		{
			name: "create signature device already exist",
			fields: fields{
				signatureDevice: map[string]domain.SignatureDeviceResponse{
					"someid": {
						ID:        "someid",
						Label:     "some label",
						Algorithm: crypto.SignatureAlgorithmRSA,
					},
				},
				deviceSignatures: make(map[string][]domain.SignatureResponse),
			},
			args: args{
				sdreq: domain.SignatureDeviceRequest{
					ID:        "someid",
					Label:     "some label",
					Algorithm: crypto.SignatureAlgorithmRSA,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &inMemorySignatureDeviceRepository{
				signatureDevice:  tt.fields.signatureDevice,
				deviceSignatures: tt.fields.deviceSignatures,
			}
			gotSdres, err := r.Create(tt.args.sdreq)
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemorySignatureDeviceRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSdres, tt.want) {
				t.Errorf("inMemorySignatureDeviceRepository.Create() = %v, want %v", gotSdres, tt.want)
			}
		})
	}
}

func Test_inMemorySignatureDeviceRepository_Get(t *testing.T) {
	type fields struct {
		signatureDevice  map[string]domain.SignatureDeviceResponse
		deviceSignatures map[string][]domain.SignatureResponse
	}
	type args struct {
		deviceId string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want domain.SignatureDeviceResponse
		wantErr   bool
	}{
		{
			name: "get device success",
			fields: fields{
				signatureDevice: map[string]domain.SignatureDeviceResponse{
					"someid": {
						ID:               "someid",
						Label:            "some label",
						Algorithm:        crypto.SignatureAlgorithmRSA,
						SignatureCounter: 2,
					},
				},
				deviceSignatures: make(map[string][]domain.SignatureResponse),
			},
			args: args{
				deviceId: "someid",
			},
			want: domain.SignatureDeviceResponse{
				ID:               "someid",
				Label:            "some label",
				Algorithm:        crypto.SignatureAlgorithmRSA,
				SignatureCounter: 2,
			},
			wantErr: false,
		},
		{
			name: "get device failure - device does not exist",
			fields: fields{
				signatureDevice:  make(map[string]domain.SignatureDeviceResponse),
				deviceSignatures: make(map[string][]domain.SignatureResponse),
			},
			args: args{
				deviceId: "someid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &inMemorySignatureDeviceRepository{
				signatureDevice:  tt.fields.signatureDevice,
				deviceSignatures: tt.fields.deviceSignatures,
			}
			gotSdres, err := r.Get(tt.args.deviceId)
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemorySignatureDeviceRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSdres, tt.want) {
				t.Errorf("inMemorySignatureDeviceRepository.Get() = %v, want %v", gotSdres, tt.want)
			}
		})
	}
}

func Test_inMemorySignatureDeviceRepository_GetAll(t *testing.T) {
	type fields struct {
		signatureDevice  map[string]domain.SignatureDeviceResponse
		deviceSignatures map[string][]domain.SignatureResponse
	}
	tests := []struct {
		name    string
		fields  fields
		want    []domain.SignatureDeviceResponse
		wantErr bool
	}{
		{
			name: "get device success - no device",
			fields: fields{
				signatureDevice:  make(map[string]domain.SignatureDeviceResponse),
				deviceSignatures: make(map[string][]domain.SignatureResponse),
			},
			want:    []domain.SignatureDeviceResponse{},
			wantErr: false,
		},
		{
			name: "get device success - with devices",
			fields: fields{
				signatureDevice: map[string]domain.SignatureDeviceResponse{
					"someid": {
						ID:               "someid",
						Label:            "some label",
						Algorithm:        crypto.SignatureAlgorithmRSA,
						SignatureCounter: 2,
					},
				},
				deviceSignatures: make(map[string][]domain.SignatureResponse),
			},
			want: []domain.SignatureDeviceResponse{{
				ID:               "someid",
				Label:            "some label",
				Algorithm:        crypto.SignatureAlgorithmRSA,
				SignatureCounter: 2,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &inMemorySignatureDeviceRepository{
				signatureDevice:  tt.fields.signatureDevice,
				deviceSignatures: tt.fields.deviceSignatures,
			}
			got, err := r.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemorySignatureDeviceRepository.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inMemorySignatureDeviceRepository.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inMemorySignatureDeviceRepository_GetAllSignature(t *testing.T) {
	type fields struct {
		signatureDevice  map[string]domain.SignatureDeviceResponse
		deviceSignatures map[string][]domain.SignatureResponse
	}
	type args struct {
		deviceId string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want []domain.SignatureResponse
		wantErr  bool
	}{
		{
			name: "get all signatures success",
			fields: fields{
				signatureDevice: make(map[string]domain.SignatureDeviceResponse),
				deviceSignatures: map[string][]domain.SignatureResponse{
					"someid": {
						{
							Signature:  "thesignature",
							SignedData: "thesigneddata",
						},
					},
				},
			},
			args: args{deviceId: "someid"},
			want: []domain.SignatureResponse{
				{
					Signature:  "thesignature",
					SignedData: "thesigneddata",
				},
			},
			wantErr: false,
		},
		{
			name: "get all signatures failure - signature device does not exist",
			fields: fields{
				signatureDevice:  make(map[string]domain.SignatureDeviceResponse),
				deviceSignatures: make(map[string][]domain.SignatureResponse),
			},
			args:    args{deviceId: "someid"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &inMemorySignatureDeviceRepository{
				signatureDevice:  tt.fields.signatureDevice,
				deviceSignatures: tt.fields.deviceSignatures,
			}
			gotSres, err := r.GetAllSignature(tt.args.deviceId)
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemorySignatureDeviceRepository.GetAllSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSres, tt.want) {
				t.Errorf("inMemorySignatureDeviceRepository.GetAllSignature() = %v, want %v", gotSres, tt.want)
			}
		})
	}
}

func Test_inMemorySignatureDeviceRepository_AddSignature(t *testing.T) {
	type fields struct {
		signatureDevice  map[string]domain.SignatureDeviceResponse
		deviceSignatures map[string][]domain.SignatureResponse
	}
	type args struct {
		deviceId string
		sres     domain.SignatureResponse
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want domain.SignatureDeviceResponse
		wantErr   bool
	}{
		{
			name: "add signature to device success",
			fields: fields{
				signatureDevice: map[string]domain.SignatureDeviceResponse{
					"someid": {
						ID:               "someid",
						Label:            "some label",
						Algorithm:        crypto.SignatureAlgorithmRSA,
						SignatureCounter: 1,
					},
				},
				deviceSignatures: make(map[string][]domain.SignatureResponse),
			},
			args: args{deviceId: "someid", sres: domain.SignatureResponse{
				Signature:  "thesignature",
				SignedData: "thesigneddata",
			}},
			want: domain.SignatureDeviceResponse{
				ID:               "someid",
				Label:            "some label",
				Algorithm:        crypto.SignatureAlgorithmRSA,
				SignatureCounter: 2,
			},
			wantErr: false,
		},
		{
			name: "add signature failure - signature device does not exist",
			fields: fields{
				signatureDevice:  make(map[string]domain.SignatureDeviceResponse),
				deviceSignatures: make(map[string][]domain.SignatureResponse),
			},
			args:    args{deviceId: "someid"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &inMemorySignatureDeviceRepository{
				signatureDevice:  tt.fields.signatureDevice,
				deviceSignatures: tt.fields.deviceSignatures,
			}
			gotSdres, err := r.AddSignature(tt.args.deviceId, tt.args.sres)
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemorySignatureDeviceRepository.AddSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSdres, tt.want) {
				t.Errorf("inMemorySignatureDeviceRepository.AddSignature() = %v, want %v", gotSdres, tt.want)
			}
		})
	}
}
