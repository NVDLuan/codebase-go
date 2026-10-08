package httpserver

import (
	"encoding/json"
	"net/http"
)

// Map lỗi domain → status nằm ở transport của từng feature (vd. order/transport/http/errors.go)
// để platform không phải import feature.

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}
