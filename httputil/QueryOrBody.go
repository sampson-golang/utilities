package httputil

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func QueryOrBody(request *http.Request, keys ...string) (map[string]string, error) {
	values := map[string]string{}
	for _, key := range keys {
		if request.URL.Query().Get(key) != "" {
			values[key] = request.URL.Query().Get(key)
		}
	}

	if request.Method == "GET" {
		return values, nil
	}

	contentType := strings.Split(request.Header.Get("Content-Type"), ";")[0]
	log.Println("Content-Type:", contentType, "|", request.Header.Get("Content-Type"))
	switch contentType {
	case "application/json":
		decoder := json.NewDecoder(request.Body)
		var data map[string]string
		err := decoder.Decode(&data)

		if err != nil {
			return values, errors.New("invalid JSON")
		}

		for _, key := range keys {
			if data[key] != "" {
				values[key] = data[key]
			}
		}
	default:
		log.Println("Form:", request.Form)

		for _, key := range keys {
			if request.FormValue(key) != "" {
				values[key] = request.Form.Get(key)
			}
		}
	}
	return values, nil
}
