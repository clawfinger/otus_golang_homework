package internalhttp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	servers "github.com/clawfinger/hw12_13_14_15_calendar/internal/server"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/storage"
)

type Handler struct {
	serverCtx *servers.ServerContext
}

func NewHandler(serverCtx *servers.ServerContext) *Handler {
	return &Handler{
		serverCtx: serverCtx,
	}
}

func UnmarshallEvent(r *http.Request) (*storage.Event, error) {
	event := &storage.Event{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(event)
	return event, err
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	event, err := UnmarshallEvent(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot unmarshall the body"))
		return
	}
	h.serverCtx.Storage.Create(event)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	event, err := UnmarshallEvent(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot unmarshall the body"))
		return
	}
	h.serverCtx.Storage.Update(event)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	event, err := UnmarshallEvent(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot unmarshall the body"))
		return
	}
	h.serverCtx.Storage.Delete(event)
}

func getTimeFromParams(r *http.Request) (time.Time, error) {
	timeString := r.URL.Query().Get("time")
	timeInt, err := strconv.ParseInt(timeString, 10, 64)
	time := time.Unix(timeInt, 0)
	return time, err
}

func (h *Handler) GetDay(w http.ResponseWriter, r *http.Request) {
	time, err := getTimeFromParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong time"))
		return
	}

	events := h.serverCtx.Storage.GetEventsForDay(time)
	body, err := json.Marshal(events)
	if err != nil {
		h.serverCtx.Logger.Info("Failed to marshall events slice to json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to marshall request result"))
	}
	w.Write(body)
}

func (h *Handler) GetWeek(w http.ResponseWriter, r *http.Request) {
	time, err := getTimeFromParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong time"))
		return
	}

	events := h.serverCtx.Storage.GetEventsForWeek(time)
	body, err := json.Marshal(events)
	if err != nil {
		h.serverCtx.Logger.Info("Failed to marshall events slice to json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to marshall request result"))
	}
	w.Write(body)
}

func (h *Handler) GetMonth(w http.ResponseWriter, r *http.Request) {
	time, err := getTimeFromParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong time"))
		return
	}

	events := h.serverCtx.Storage.GetEventsForMonth(time)
	body, err := json.Marshal(events)
	if err != nil {
		h.serverCtx.Logger.Info("Failed to marshall events slice to json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to marshall request result"))
	}
	w.Write(body)
}
