package internalhttp

import (
	"context"
	"net/http"
)

type Server struct { // TODO
	server    *http.Server
	mux       *http.ServeMux
	handler   Handler
	serverCtx *ServerContext
}

func NewServer(serverCtx *ServerContext) *Server {
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
