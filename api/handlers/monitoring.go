package handlers

import (
	"net/http"
)

type HealthCheckResult struct {
	Status string `json:"status"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	SendJsonResponse(w, HealthCheckResult{Status: "ok"}, http.StatusOK)
}
