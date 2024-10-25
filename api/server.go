package api

import (
	"encoding/json"
	"net/http"

	"github.com/GiacomoCortesi/gosign/domain"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress          string
	signatureDeviceService domain.SignatureDeviceService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, service domain.SignatureDeviceService) *Server {
	return &Server{
		listenAddress:          listenAddress,
		signatureDeviceService: service,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))

	mux.Handle("/api/v0/devices", http.HandlerFunc(s.SignatureDevicesHandler))
	mux.Handle("/api/v0/devices/{id}", http.HandlerFunc(s.SignatureDeviceHandler))
	mux.Handle("/api/v0/devices/{id}/signatures", http.HandlerFunc(s.SignTransactionHandler))
	return http.ListenAndServe(s.listenAddress, corsMiddleware(mux))
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	// set appropriate content type header
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	// set appropriate content type header
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
