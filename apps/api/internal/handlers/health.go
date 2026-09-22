package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}

// Health reports whether the service and its database are reachable.
// Used by the frontend's /health page and by container health checks.
func Health(conn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{Status: "ok", DB: "ok"}
		status := http.StatusOK

		if err := conn.PingContext(r.Context()); err != nil {
			resp.Status = "degraded"
			resp.DB = "unreachable"
			status = http.StatusServiceUnavailable
		}

		writeJSON(w, status, resp)
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
