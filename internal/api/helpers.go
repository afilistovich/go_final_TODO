package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func parseID(r *http.Request) (int64, error) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return 0, fmt.Errorf("id parameter is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id format")
	}
	return id, nil
}
