package monitoring

import (
	"net/http"

	"github.com/andriystech/lgc/api/handlers/common"
)

type HealthCheckResult struct {
	Status string `json:"status"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	common.SendJsonResponse(w, HealthCheckResult{Status: "ok"}, http.StatusOK)
}
