package proxyserver

import (
	"bytes"
	"io"
	"net/http"
)

// Routes ...
func (s *Server) Routes() {
	s.mux.HandleFunc("/", s.CacheMiddleware(s.ProxyHandler()))
}

// ProxyHandler ...
func (s *Server) ProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		urn := s.config.URL + r.URL.String()

		res, err := http.Get("http://" + urn)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			s.logger.Error("Error reading body: ", err)
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}
		defer res.Body.Close()

		s.cache.CacheData(urn, body)

		_, err = io.Copy(w, bytes.NewReader(body))
		if err != nil {
			s.logger.Error(err)
			http.Error(w, "500 - internal server error", http.StatusInternalServerError)
		}

	}
}

// CacheMiddleware ...
func (s *Server) CacheMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		urn := s.config.URL + r.URL.String()
		s.logger.Info("making request to: ", "http://"+urn)

		cacheReader, err := s.cache.GetCacheReader(urn)
		if err != nil {
			s.logger.Error("cacheReader error: ", err)
			http.Error(w, "500 - internal server error", http.StatusInternalServerError)
			return
		}

		if cacheReader != nil {
			defer cacheReader.Close()
			s.logger.Debug("Берем данные из кэша.")
			// Have cache data
			_, err = io.Copy(w, cacheReader)
			if err != nil {
				s.logger.Error(err)
			}
			return
		}

		s.logger.Debug("Кэш не найден. Делаем запрос.")
		f(w, r)
	}
}
