package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/antoniofmoliveira/fullcycle-desafio-ratelimiter/internal/config"
)

type Server struct {
	Server     *http.Server
	Config     config.Config
	Handles    []Handle
	Middleware func(http.Handler) http.Handler
}

// NewServer returns a new Server instance with the given configuration.
// It creates an *http.Server instance with the given port,
// sets the server's middleware to the result of GetAllMiddlewares
// and sets the server's Handles to AvailableHandles.
func NewServer(c config.Config) Server {

	s := &Server{
		Server: &http.Server{
			Addr: fmt.Sprintf(":%d", c.Port),
		},
		Config:     c,
		Middleware: GetAllMiddlewares(&c),
		Handles:    AvailableHandles,
	}

	return *s
}

// AddHandle adds a new handle to the server's list of handles. It takes
// a Handle as argument and appends it to the server's Handles slice.
func (s *Server) AddHandle(h Handle) {
	s.Handles = append(s.Handles, h)
}

// Start initializes the HTTP server by registering all the available
// handles with their respective paths and middleware. It then starts
// listening for incoming requests on the configured port, logging the
// server start information. The function blocks and serves requests
// until the server is shut down.
func (s *Server) Start() {
	for _, h := range s.Handles {
		http.Handle(h.Path, s.Middleware(h.Handler))
	}
	slog.Info("Server started", "port", s.Config.Port)
	s.Server.ListenAndServe()
}
