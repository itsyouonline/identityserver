package goraml

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// APIError represents any error condition that returned by API server.
// i.e. non 2xx status code
type APIError struct {
	Code       int         // http status code
	RawMessage []byte      // raw response body
	Message    interface{} // parsed response body, if possible, it might be nil
}

// Error implements error interface
func (e APIError) Error() string {
	return string(e.RawMessage)
}

// NewAPIError creates new APIError object.
// If data is nil, it raw message won't be parsed.
func NewAPIError(resp *http.Response, data interface{}) error {
	// read response body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// creates object
	goramlErr := APIError{
		Code:       resp.StatusCode,
		RawMessage: b,
	}

	// no need to parse in case of `data` is nil
	if data == nil {
		return goramlErr
	}

	// parse based on content type
	switch resp.Header.Get("Content-Type") {
	case "application/json":
		err = json.Unmarshal(b, data)
	case "text/plain":
		_, err = fmt.Sscanf(string(b), "%v", data)
	default:
		return goramlErr
	}

	if err == nil {
		goramlErr.Message = data
	}
	return goramlErr
}
