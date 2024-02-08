package main

import (
	inst "broker-service/proto_files/institution_proto"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

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
// Modified rpcErrorJson function to handle gRPC status codes
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

    return app.writeJSON(w, statusCode, jsonResponse{
        Error:   true,
        Message: st.Message(),
    })
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

func (app *Config) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *Config) readListOfIntValues(qs url.Values, key string, defaultValue []int64) ([]int64, error) {
	s := qs.Get(key)

	if s == "" {
		return defaultValue, nil
	}
	res := []int64{}
	for _, v := range strings.Split(s, ",") {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return defaultValue, err
		}
		res = append(res, i)
	}

	return res, nil
}

func (app *Config) readInt(qs url.Values, key string, defaultValue int) (int, error) {
	s := qs.Get(key)
	if s == "" {
		return defaultValue, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (app *Config) getWorkHours(requestPayload RequestPayload) []*inst.WorkingHours {
	var workHours []*inst.WorkingHours
	for _, wh := range requestPayload.Institution.WorkingHours {
		temp := &inst.WorkingHours{
			Day:   int32(wh.Day),
			Open:  wh.Open,
			Close: wh.Close,
		}
		workHours = append(workHours, temp)
	}
	return workHours
}
