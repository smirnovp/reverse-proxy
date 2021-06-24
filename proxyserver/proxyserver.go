package proxyserver

import (
	"context"
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
func New(config *Config, cache *cache.Storage) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		mux:    http.NewServeMux(),
		cache:  cache,
	}
}

// Start ...
func (s *Server) Start() error {
	if err := s.configureProxyServer(); err != nil {
		return err
	}
	if err := s.configureLogger(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := s.cache.Run(ctx); err != nil {
		return err
	}

	s.logger.Infof("Server is starting at port%s ...", s.config.Port)

	s.Routes()
	if err := http.ListenAndServe(s.config.Port, s.mux); err != nil {
		return err
	}

	return nil
}

func (s *Server) configureProxyServer() error {
	return nil
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}
