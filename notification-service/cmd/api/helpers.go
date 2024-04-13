package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

var (
	ErrNotFound       = errors.New("resource not found")
	ErrBadRequest     = errors.New("bad request")
	ErrInternal       = errors.New("data server error")
	ErrAuthentication = errors.New("authentication failed")
)

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// func (app *Config) readStringParam(r *http.Request, key string) (string, error) {
// 	value := chi.URLParam(r, key)
// 	if value == "" {
// 		return "", ErrBadRequest
// 	}

// 	return value, nil
// }

func (app *Config) readIntParam(r *http.Request, key string) (int64, error) {
	value := chi.URLParam(r, key)
	if value == "" {
		return 0, ErrBadRequest
	}
	//parse string to int64
	return strconv.ParseInt(value, 10, 64)
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}
func (app *Config) errorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

func (app *Config) rpcErrorJson(w http.ResponseWriter, err error) error {
	// Extract the gRPC status from the error
	st, ok := status.FromError(err)
	if !ok {
		// This is not a gRPC error
		return app.writeJSON(w, http.StatusInternalServerError, jsonResponse{
			Error:   true,
			Message: "An unexpected error occurred",
		})
	}

	// Map gRPC status codes to HTTP status codes
	var statusCode int
	switch st.Code() {
	case codes.InvalidArgument:
		statusCode = http.StatusBadRequest
	case codes.NotFound:
		statusCode = http.StatusNotFound
	case codes.AlreadyExists:
		statusCode = http.StatusConflict
	case codes.PermissionDenied:
		statusCode = http.StatusForbidden
	case codes.Unauthenticated:
		statusCode = http.StatusUnauthorized
	default:
		statusCode = http.StatusInternalServerError
	}

	// Simplify the message by removing the gRPC error prefix if present
	errorMsg := strings.TrimPrefix(st.Message(), "rpc error: code = "+st.Code().String()+" desc = ")

	return app.writeJSON(w, statusCode, jsonResponse{
		Error:   true,
		Message: errorMsg,
	})
}
