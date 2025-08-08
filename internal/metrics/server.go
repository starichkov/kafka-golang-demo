// Package metrics provides HTTP server for exposing Prometheus metrics.
package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"kafka-golang-demo/internal/logging"
)

// Server provides HTTP server for Prometheus metrics endpoint
type Server struct {
	server *http.Server
	addr   string
}

// NewServer creates a new metrics HTTP server
func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	
	// Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		server: server,
		addr:   addr,
	}
}

// Start starts the metrics HTTP server
func (s *Server) Start() error {
	logging.Logger.Info("Starting metrics server", "address", s.addr)
	
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully stops the metrics HTTP server
func (s *Server) Stop(ctx context.Context) error {
	logging.Logger.Info("Stopping metrics server")
	return s.server.Shutdown(ctx)
}

// Address returns the server address
func (s *Server) Address() string {
	return s.addr
}