package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/antoniofmoliveira/fullcycle-desafio-ratelimiter/internal/config"
)

type server struct {
	Server     *http.Server
	handles    []handle
	middleware func(http.Handler) http.Handler
}

// NewServer returns a new Server instance with the given configuration.
// It creates an *http.Server instance with the given port,
// sets the server's middleware to the result of GetAllMiddlewares
// and sets the server's Handles to AvailableHandles.
func NewServer(c config.Config) *server {
	return &server{
		Server:     &http.Server{Addr: fmt.Sprintf(":%d", c.Port)},
		middleware: getAllMiddlewares(&c),
		handles:    availableHandles,
	}
}

// AddHandle adds a new handle to the server's list of handles. It takes
// a Handle as argument and appends it to the server's Handles slice.
func (s *server) AddHandle(h handle) {
	s.handles = append(s.handles, h)
}

// Start initializes the HTTP server by registering all the available
// handles with their respective paths and middleware. It then starts
// listening for incoming requests on the configured port, logging the
// server start information. The function blocks and serves requests
// until the server is shut down.
func (s *server) Start() {
	for _, h := range s.handles {
		http.Handle(h.path, s.middleware(h.handler))
	}
	slog.Info("Server started", "Address", s.Server.Addr)
	s.Server.ListenAndServe()
}
