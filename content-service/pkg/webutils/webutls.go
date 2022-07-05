package webutils

import (
	"encoding/json"
	"net/http"
)

// WriteHTTPCode gets corresponding status text for err status code
// and writes both header and status text with w
func WriteHTTPCode(w http.ResponseWriter, err int) {
	w.WriteHeader(err)
	w.Write([]byte(http.StatusText(err)))
	return
}

// WriteHTTPCodeJSON gets corresponding status text for err status code
// adds "status" field to messages map with status text, json marshals the map
// and writes the response using w. If messages is nil, a new map is created
func WriteHTTPCodeJSON(w http.ResponseWriter, err int, messages map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err)
	if messages == nil {
		messages = make(map[string]string)
	}
	messages["status"] = http.StatusText(err)
	result, _ := json.Marshal(messages)
	w.Write(result)
}
