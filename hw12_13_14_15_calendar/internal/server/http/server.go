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
	handler := NewHandler(serverCtx.Logger)
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", loggingMiddleware(handler.helloWorld, serverCtx.Logger))
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
