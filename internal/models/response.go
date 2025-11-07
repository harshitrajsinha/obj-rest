package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Response defines the strucutre of response data that will be sent for a request
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendResponse constructs the response to be sent for the request
func SendResponse(w http.ResponseWriter, code int, message string, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")

	responseToSend := &Response{
		Status:  http.StatusText(code),
		Message: message,
		Data:    data,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(responseToSend); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "Internal Server Error", "message": "Please try again later"}`))
		return fmt.Errorf("error sending response, %w", err)
	}

	w.WriteHeader(code)
	buf.WriteTo(w)

	return nil

}
