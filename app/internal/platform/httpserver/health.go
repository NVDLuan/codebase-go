package httpserver

import (
	"context"
	"net/http"
	"time"
)

// Pinger là phụ thuộc cần sẵn sàng (vd. *pgxpool.Pool). nil = không kiểm tra.
type Pinger interface {
	Ping(ctx context.Context) error
}

// RegisterHealth: liveness không kiểm tra phụ thuộc; readiness ping DB.
func RegisterHealth(mux *http.ServeMux, deps ...Pinger) {
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, _ *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		for _, d := range deps {
			if d == nil {
				continue
			}
			if err := d.Ping(ctx); err != nil {
				WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
				return
			}
		}
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
}
