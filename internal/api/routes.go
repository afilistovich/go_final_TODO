package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/api/nextdate", nextDayHandler)
}
