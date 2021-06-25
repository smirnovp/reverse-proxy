package proxyserver

import (
	"context"
	"fmt"
	"net/http"
	"reverse-proxy/cache"

	"github.com/sirupsen/logrus"
)

// Server ...
type Server struct {
	config *Config
	logger *logrus.Logger
	mux    *http.ServeMux
	cache  *cache.Storage
}

// New ...
func New(logger *logrus.Logger, config *Config, cache *cache.Storage) *Server {
	return &Server{
		config: config,
		logger: logger,
		mux:    http.NewServeMux(),
		cache:  cache,
	}
}

// Start ...
func (s *Server) Start() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := s.cache.Run(ctx); err != nil {
		return err
	}

	s.Routes()

	fmt.Printf("Server is starting at port%s ...\n", s.config.Port)

	if err := http.ListenAndServe(s.config.Port, s.mux); err != nil {
		return err
	}

	return nil
}
