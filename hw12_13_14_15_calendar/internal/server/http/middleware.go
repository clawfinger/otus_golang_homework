package internalhttp

import (
	"net/http"
	"time"

	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(next http.HandlerFunc, logger logger.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workStart := time.Now()
		next(w, r)
		logger.Info("Request", "Ip:", r.RemoteAddr, "Duration:",
			time.Since(workStart), "Method:", r.Method, "Path", r.URL.Path)
	})
}
