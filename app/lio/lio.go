package lio

import (
	"encoding/json"
	"net/http"
)

// HandleGETResponse : adds data to the body of a GET response
func HandleGETResponse(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

// HandlePOSTResponse : returns a response to a POST request
func HandlePOSTResponse(w http.ResponseWriter) {
	json.NewEncoder(w).Encode("OK!")
}

// DecodePOSTBody : decodes the POST body into an object (v)
func DecodePOSTBody(r *http.Request, v interface{}) error {
	body := r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&v)
	return err
}
