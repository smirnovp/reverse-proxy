package proxyserver

import (
	"io"
	"net/http"
)

// Routes ...
func (s *Server) Routes() {
	s.mux.HandleFunc("/", s.ProxyHandler())
}

// ProxyHandler ...
func (s *Server) ProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("making request to: ", "http://"+s.config.URL+r.URL.String())

		res, err := http.Get("http://" + s.config.URL + r.URL.String())
		if err != nil {
			s.logger.Error(err)
		}

		_, err = io.Copy(w, res.Body)
		if err != nil {
			s.logger.Error(err)
		}

		defer res.Body.Close()
	}
}
