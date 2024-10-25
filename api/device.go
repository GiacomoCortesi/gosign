/*
Package API provides a REST API interface to the gosign microservice.

It registers API routes and expose HTTP handlers to manage signature devices and to sign data
with the signature devices.

It include CORS middleware for localhost for development purposes.

Check openapi.yaml definition for more in depth REST API documentation.
*/
package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GiacomoCortesi/gosign/domain"
)

// SignatureDevicesHandler dispatch signature devices requests
func (s *Server) SignatureDevicesHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		s.GetAllSignatureDevice(response, request)
	case http.MethodPost:
		s.CreateSignatureDevice(response, request)
	default:
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
	}
}

// GetAllSignatureDevice fetch available signature devices
func (s *Server) GetAllSignatureDevice(response http.ResponseWriter, request *http.Request) {
	sdres, err := s.signatureDeviceService.GetAll()
	if err != nil {
		WriteErrorResponse(response, http.StatusServiceUnavailable, []string{
			http.StatusText(http.StatusServiceUnavailable),
		})
		return
	}
	WriteAPIResponse(response, http.StatusOK, sdres)
}

// CreateSignatureDevice create a new signature device
func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	var sdreq domain.SignatureDeviceRequest
	if err := json.NewDecoder(request.Body).Decode(&sdreq); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
		})
		return
	}

	sdres, err := s.signatureDeviceService.Create(sdreq)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSignatureDeviceAlreadyExist):
			WriteErrorResponse(response, http.StatusConflict, []string{
				http.StatusText(http.StatusConflict),
			})
		default:
			WriteErrorResponse(response, http.StatusServiceUnavailable, []string{
				http.StatusText(http.StatusServiceUnavailable),
			})
		}
		return
	}

	WriteAPIResponse(response, http.StatusCreated, sdres)
}

// SignatureDeviceHandler dispatch signature device requests
func (s *Server) SignatureDeviceHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		s.GetSignatureDevice(response, request)
	default:
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
	}
}

// GetSignatureDevice fetch a signature device given its ID
func (s *Server) GetSignatureDevice(response http.ResponseWriter, request *http.Request) {
	deviceId := request.PathValue("id")
	sdres, err := s.signatureDeviceService.Get(deviceId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSignatureDeviceNotFound):
			WriteErrorResponse(response, http.StatusNotFound, []string{
				http.StatusText(http.StatusNotFound),
			})
		default:
			WriteErrorResponse(response, http.StatusServiceUnavailable, []string{
				http.StatusText(http.StatusServiceUnavailable),
			})
		}
		return
	}
	WriteAPIResponse(response, http.StatusOK, sdres)
}

// SignTransactionHandler dispatch transaction signature requests
func (s *Server) SignTransactionHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		s.SignTransaction(response, request)
	case http.MethodGet:
		s.GetDeviceSignatures(response, request)
	default:
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
	}
}

// SignTransaction sign the request data using the appropriate signature device
func (s *Server) SignTransaction(response http.ResponseWriter, request *http.Request) {
	deviceId := request.PathValue("id")

	var sreq domain.SignatureRequest
	if err := json.NewDecoder(request.Body).Decode(&sreq); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
		})
		return
	}

	sres, err := s.signatureDeviceService.SignTransaction(deviceId, sreq.Data)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSignatureDeviceNotFound):
			WriteErrorResponse(response, http.StatusNotFound, []string{
				http.StatusText(http.StatusNotFound),
			})
		default:
			WriteErrorResponse(response, http.StatusServiceUnavailable, []string{
				http.StatusText(http.StatusServiceUnavailable),
			})
		}
		return
	}
	WriteAPIResponse(response, http.StatusOK, sres)
}

// GetDeviceSignatures fetch all transaction signatures for the specified signature device
func (s *Server) GetDeviceSignatures(response http.ResponseWriter, request *http.Request) {
	deviceId := request.PathValue("id")

	sres, err := s.signatureDeviceService.GetAllSignature(deviceId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSignatureDeviceNotFound):
			WriteErrorResponse(response, http.StatusNotFound, []string{
				http.StatusText(http.StatusNotFound),
			})
		default:
			WriteErrorResponse(response, http.StatusServiceUnavailable, []string{
				http.StatusText(http.StatusServiceUnavailable),
			})
		}
		return
	}
	WriteAPIResponse(response, http.StatusOK, sres)
}
