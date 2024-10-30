package sat

import (
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
)

const (
	// SUCCESS_OK contains a success message
	SUCCESS_OK = "OK"
	// INVALID_SIGNATURE contains invalid signature message
	INVALID_SIGNATURE = "INVALID_SIGNATURE"
	// INVALID_PAYLOAD contains invalid payload message
	INVALID_PAYLOAD = "INVALID_PAYLOAD"

	// EMPTY_CLIENT_ID contains an empty client id error message
	EMPTY_CLIENT_ID = "client id can't be empty"
	// EMPTY_CLIENT_SECRET contains an empty client secret error message
	EMPTY_CLIENT_SECRET = "client secret can't be empty"
	// EMPTY_CLIENT_PRIVATE_KEY contains an empty client private key error message
	EMPTY_CLIENT_PRIVATE_KEY = "client private key can't be empty"
)

// APIResponseError contains interface error method
type APIResponseError interface {
	Error() string
	Code() string
	Status() string
	Detail() string
}

// ErrorResponse wrapper api error
type ErrorResponse struct {
	Errors []*ErrorObject `json:"errors"`
}

// ErrorObject is jsonapi.ErrorObject
type ErrorObject jsonapi.ErrorObject

// Error will convert all detail errors to one string
func (e *ErrorResponse) Error() string {
	if len(e.Errors) <= 0 {
		return ""
	}

	return fmt.Sprintf("%s - %s - %s\n", e.Errors[0].Status, e.Errors[0].Code, e.Errors[0].Detail)
}

// Code will parse code error and return it
func (e *ErrorResponse) Code() string {
	if len(e.Errors) <= 0 {
		return ""
	}

	return e.Errors[0].Code
}

// Status will parse status error and return it
func (e *ErrorResponse) Status() string {
	if len(e.Errors) <= 0 {
		return ""
	}

	return e.Errors[0].Status
}

// Detail will parse the detail error and return it
func (e *ErrorResponse) Detail() string {
	if len(e.Errors) <= 0 {
		return ""
	}

	return e.Errors[0].Detail
}

// APIInternalError for internal error produces by non-SAT server
type APIInternalError interface {
	Error() string
	Response() *http.Response
}

// InternalError wrapper internal error http response
type InternalError struct {
	resp *http.Response
}

// Error will return http status code and http status as string
func (i *InternalError) Error() string {
	return fmt.Sprintf("%d - %s\n", i.resp.StatusCode, i.resp.Status)
}

// Response will return http raw response
func (i *InternalError) Response() *http.Response {
	return i.resp
}
