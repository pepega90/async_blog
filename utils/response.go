package utils

import (
	"encoding/json"
	"net/http"
)

func ResponsWithError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

func RespondJSON(w http.ResponseWriter, code int, payload any) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
