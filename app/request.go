package app

import "net/http"

func FormValueOrDefault(request *http.Request, key string, defaultValue string) string {
	value := request.FormValue(key)
	if value != "" {
		return value
	}
	return defaultValue
}
