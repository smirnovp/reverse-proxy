package proxyserver

import "github.com/sirupsen/logrus"

// Server ...
type Server struct {
	config *Config
	logger *logrus.Logger
}

// New ...
func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
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
	s.logger.Info("Server is starting...")
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
