package server

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Server struct
type Server struct {
	port   int
	router *http.ServeMux
}

// NewServer creates a Server instance with mux
func NewServer(port int) *Server {
	return &Server{
		port:   port,
		router: http.NewServeMux(),
	}
}

// AddRoute adds a route with a handler
func (s *Server) AddRoute(path string, handler http.Handler) {
	s.router.Handle(path, handler)
}

// Run starts the http server
func (s *Server) Run() {
	log.Info(fmt.Sprintf("Starting commit-message-bot on port %d", s.port))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router))
}
