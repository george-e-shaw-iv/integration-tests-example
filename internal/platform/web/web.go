package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ErrInternalServer is a generic internal server error.
var ErrInternalServer = errors.New("internal server error")

// Response is the format used for all the responses.
type Response struct {
	Results interface{}     `json:"results"`
	Errors  []ResponseError `json:"errors,omitempty"`
}

// ResponseError is the format used for response errors.
type ResponseError struct {
	Message string `json:"message"`
}

// Error implements the error interface
func (a ResponseError) Error() string {
	return a.Message
}

// Respond sends a response with a status code.
func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}, errs ...error) {
	var respErrs []ResponseError

	if len(errs) > 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("error while serving request")

			respErrs = append(respErrs, ResponseError{Message: err.Error()})
		}
	}

	resp := Response{
		Results: data,
		Errors:  respErrs,
	}

	writeResponse(w, r, code, &resp)
}

// RespondError sends an error response with a status code. The error is automatically logged for you.
// If the error implements StatusCoder, the provided status code will be used.
func RespondError(w http.ResponseWriter, r *http.Request, code int, err error) {
	log.WithFields(log.Fields{
		"error": err,
	}).Error("error while serving request")

	if code >= http.StatusInternalServerError && code != http.StatusServiceUnavailable && code != http.StatusNotImplemented {

		// Respond with generic error. Error messages and and codes may potentially contain
		// sensitive information or help an attacker.
		code = http.StatusInternalServerError
		//err = ErrInternalServer
	}

	resp := Response{
		Errors: []ResponseError{
			{
				Message: err.Error(),
			},
		},
	}

	writeResponse(w, r, code, &resp)
}

// writeResponse marshals the response to json and writes it to the response writer.
func writeResponse(w http.ResponseWriter, r *http.Request, code int, resp *Response) {
	if code == http.StatusNoContent || resp == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}
