package internalhttp

import (
	"net/http"

	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
)

type Handler struct {
	log logger.Logger
}

func NewHandler(logger logger.Logger) *Handler {
	return &Handler{
		log: logger,
	}
}

func (h *Handler) helloWorld(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Hello world")
	w.Write([]byte("Hello world"))
}
