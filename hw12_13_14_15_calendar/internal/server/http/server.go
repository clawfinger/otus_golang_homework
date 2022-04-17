package internalhttp

import (
	"context"
	"net/http"

	servers "github.com/clawfinger/hw12_13_14_15_calendar/internal/server"
)

type Server struct { // TODO
	server    *http.Server
	mux       *http.ServeMux
	handler   Handler
	serverCtx *servers.ServerContext
}

func NewServer(serverCtx *servers.ServerContext) *Server {
	handler := NewHandler(serverCtx)
	mux := http.NewServeMux()
	mux.HandleFunc("/create", loggingMiddleware(handler.Create, serverCtx.Logger))
	mux.HandleFunc("/update", loggingMiddleware(handler.Update, serverCtx.Logger))
	mux.HandleFunc("/delete", loggingMiddleware(handler.Delete, serverCtx.Logger))
	mux.HandleFunc("/getday", loggingMiddleware(handler.GetDay, serverCtx.Logger))
	mux.HandleFunc("/getweek", loggingMiddleware(handler.GetWeek, serverCtx.Logger))
	mux.HandleFunc("/getmonth", loggingMiddleware(handler.GetMonth, serverCtx.Logger))

	server := &http.Server{
		Addr:    serverCtx.Cfg.Data.HTTP.Addr,
		Handler: mux,
	}

	return &Server{
		server:    server,
		mux:       mux,
		handler:   *handler,
		serverCtx: serverCtx,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.server.ListenAndServe()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.Shutdown(ctx)
	return nil
}

// TODO
