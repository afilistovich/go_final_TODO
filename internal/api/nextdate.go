package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/afilistovich/go_final_TODO/internal/calc"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.URL.Query().Get("now")
	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(calc.DateLayout, nowStr)
		if err != nil {
			writeError(w, "invalid 'now' parameter format, expected YYYYMMDD", http.StatusBadRequest)
			return
		}
	}
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	result, err := calc.NextDate(now, date, repeat)
	if err != nil {
		slog.Warn("NextDate calculation failed", "error", err, "date", date, "repeat", repeat)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(result))
}
