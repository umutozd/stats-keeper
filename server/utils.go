package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/umutozd/stats-keeper/storage"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalRequestBody reads the request body, unmarshals it to given object and
// returns the unmarshaled object. If it fails, unmarshalRequestBody also handles
// sending the appropriate response body and status code.
func unmarshalRequestBody[T proto.Message](w http.ResponseWriter, r *http.Request, unmarshalTo T) T {
	var nilResult T // result to return when nil is intended to be returned
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid http request body", err)
		return nilResult
	}
	if err := protojson.Unmarshal(body, unmarshalTo); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid json request body", err)
		return nilResult
	}
	return unmarshalTo
}

type apiError struct {
	Message string `json:"message"`
	Err     error  `json:"error,omitempty"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("apiError: %s; %v", e.Message, e.Err)
}

func writeStorageError(w http.ResponseWriter, err error) {
	statusCode, message, wrappedErr := storage.ToHttpError(err)
	writeErrorResponse(w, statusCode, message, wrappedErr)
}

// writeErrorResponse writes an apiError to w with statusCode.
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	writeJsonResponse(w, statusCode, &apiError{Message: message, Err: err})
}

// writeJsonResponse json-marshals the given data and writes it along with statusCode
// to w. The type of marshaler used is based upon the data's type with default one being
// encoding/json.
func writeJsonResponse(w http.ResponseWriter, statusCode int, data any) {
	var resp []byte
	var err error
	if protoMsg, ok := data.(protoreflect.ProtoMessage); ok {
		// marshal using protojson package
		resp, err = protojson.Marshal(protoMsg)
	} else if b, ok := data.([]byte); ok {
		// data is already marshaled
		resp = b
	} else {
		// marshal using builtin json package
		resp, err = json.Marshal(data)
	}

	if err != nil {
		// fall back to pre-defined error message
		logrus.WithError(err).Error("writeJsonResponse: error marshaling response")
		resp = []byte(fmt.Sprintf(`{"message":"unable to encode http response","error":"%v"}`, err))
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	if _, err = w.Write(resp); err != nil {
		logrus.WithError(err).Error("writeJsonResponse: error writing http response")
	}
}

// validateRequestMethod checks if given request's method is in the given allowed methods. If so, it returns true.
// Otherwise, it sets 405 status and Allow header with given allowed methods and returns false.
func validateRequestMethod(w http.ResponseWriter, r *http.Request, allowedMethods ...string) (isValid bool) {
	for _, m := range allowedMethods {
		if r.Method == m {
			return true
		}
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
	return false
}
