package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiVersion = "/api/v1"
)

// JSONError as rest response
type JSONError struct {
	Error string `json:"error"`
}

// VOne prefix endpoints with api version 1
func VOne(path string) string {
	switch {
	case len(path) > 0:
		return fmt.Sprintf("%s/%s", apiVersion, path)
	default:
		return apiVersion
	}
}

// RenderJSON as rest response
func RenderJSON(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// log.Println("Render response", obj)
	return json.NewEncoder(w).Encode(obj)
}

// ParseJSON from incoming requests
func ParseJSON(req *http.Request, obj interface{}) error {

	if err := json.NewDecoder(req.Body).Decode(obj); err != nil {
		return fmt.Errorf("Invalid request json. Error:%s", err)
	}

	// log.Println("Parsed request ", obj)
	return nil
}
