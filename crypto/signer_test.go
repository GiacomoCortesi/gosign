package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"reflect"
	"testing"
)

func TestCreateSigner(t *testing.T) {
	rsaGgen := RSAGenerator{}
	rsaKp, err := rsaGgen.Generate()
	if err != nil {
		t.Fatalf("test setup failed, cannot create key, error: %s", err)
	}
	_, rsaPrivate, err := NewRSAMarshaler().Marshal(*rsaKp)
	if err != nil {
		t.Fatalf("test setup failed, cannot marshal key, error: %s", err)
	}

	eccGen := ECCGenerator{}
	eccKp, err := eccGen.Generate()
	if err != nil {
		t.Fatalf("test setup failed, cannot create key, error: %s", err)
	}
	_, eccPrivate, err := NewECCMarshaler().Encode(*eccKp)
	if err != nil {
		t.Fatalf("test setup failed, cannot marshal key, error: %s", err)
	}

	type args struct {
		a  SignatureAlgorithm
		pk []byte
	}
	tests := []struct {
		name    string
		args    args
		wantS   Signer
		wantErr bool
	}{
		{
			name:    "invalid algorithm",
			args:    args{SignatureAlgorithm(5), rsaPrivate},
			wantErr: true,
			wantS:   nil,
		},
		{
			name:    "ECC algorithm",
			args:    args{SignatureAlgorithmECC, eccPrivate},
			wantErr: false,
			wantS:   &ECCSigner{},
		},
		{
			name:    "RSA algorithm",
			args:    args{SignatureAlgorithmRSA, rsaPrivate},
			wantErr: false,
			wantS:   &RSASigner{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := NewSignerFactory().CreateSigner(tt.args.a, tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSigner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(gotS) != reflect.TypeOf(tt.wantS) {
				t.Errorf("CreateSigner() = %T, want %T", gotS, tt.wantS)
			}
		})
	}
}

func TestRSASigner_Sign(t *testing.T) {
	gen := RSAGenerator{}
	kp, err := gen.Generate()
	pk := kp.Private
	if err != nil {
		t.Fatalf("test setup failed, cannot create key, error: %s", err)
	}

	type fields struct {
		pk *rsa.PrivateKey
	}
	type args struct {
		dataToBeSigned []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "sign empty data",
			fields:  fields{pk},
			args:    args{[]byte("")},
			wantErr: false,
		},
		{
			name:    "sign valid data",
			fields:  fields{pk},
			args:    args{[]byte("some-valid-data-to-sign")},
			wantErr: false,
		},
		{
			name:    "sign data with special characters",
			fields:  fields{pk},
			args:    args{[]byte("!\"#$%&'()*+,-./:;<=>?@[]^_`{|}~")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RSASigner{
				pk: tt.fields.pk,
			}
			got, err := s.Sign(tt.args.dataToBeSigned)
			if (err != nil) != tt.wantErr {
				t.Errorf("RSASigner.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			hashed := sha256.Sum256(tt.args.dataToBeSigned)
			err = rsa.VerifyPKCS1v15(kp.Public, crypto.SHA256, hashed[:], got)
			if err != nil {
				t.Errorf("RSASigner.Sign() signed data fails verification")
			}
		})
	}
}

func TestECCSigner_Sign(t *testing.T) {
	gen := ECCGenerator{}
	kp, err := gen.Generate()
	pk := kp.Private
	if err != nil {
		t.Fatalf("test setup failed, cannot create key, error: %s", err)
	}

	type fields struct {
		pk *ecdsa.PrivateKey
	}
	type args struct {
		dataToBeSigned []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "sign empty data",
			fields:  fields{pk},
			args:    args{[]byte("")},
			wantErr: false,
		},
		{
			name:    "sign valid data",
			fields:  fields{pk},
			args:    args{[]byte("some-valid-data-to-sign")},
			wantErr: false,
		},
		{
			name:    "sign data with special characters",
			fields:  fields{pk},
			args:    args{[]byte("!\"#$%&'()*+,-./:;<=>?@[]^_`{|}~")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ECCSigner{
				pk: tt.fields.pk,
			}
			got, err := s.Sign(tt.args.dataToBeSigned)
			if (err != nil) != tt.wantErr {
				t.Errorf("ECCSigner.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			hashed := sha256.Sum256(tt.args.dataToBeSigned)
			ok := ecdsa.VerifyASN1(kp.Public, hashed[:], got)
			if !ok {
				t.Errorf("ECCSigner.Sign() signed data fails verification")
			}
		})
	}
}
