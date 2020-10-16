package lio

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
)

// HandleGETResponse : adds data to the body of a GET response
func HandleGETResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
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

// SetCookie :
func SetCookie(encoder *securecookie.SecureCookie, w http.ResponseWriter, key string, val string) {
	value := map[string]string{
		key: val,
	}
	if encoded, err := encoder.Encode("cookie-name", value); err == nil {
		cookie := http.Cookie{
			Name:     "cookie-name",
			Path:     "/",
			Value:    encoded,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	}
}

// ReadCookie :
func ReadCookie(encoder *securecookie.SecureCookie, r *http.Request, key string) string {
	if cookie, err := r.Cookie("cookie-name"); err == nil {
		value := make(map[string]string)
		if err = encoder.Decode("cookie-name", cookie.Value, &value); err == nil {
			return value[key]
		}
	}
	return ""
}

// GetIntParam :
func GetIntParam(ps httprouter.Params, name string) int {
	value, _ := strconv.Atoi(ps.ByName(name))
	return value
}
