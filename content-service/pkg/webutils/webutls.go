package webutils

import (
	"encoding/json"
	"net/http"
)

func WriteHTTPCode(w http.ResponseWriter, err int) {
	w.WriteHeader(err)
	w.Write([]byte(http.StatusText(err)))
	return
}

func WriteHTTPCodeJSON(w http.ResponseWriter, err int, messages map[string]string) {
	w.WriteHeader(err)
	if messages == nil {
		messages = make(map[string]string)
	}
	messages["status"] = http.StatusText(err)
	result, _ := json.Marshal(messages)
	w.Write(result)
}
