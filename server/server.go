package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umutozd/stats-keeper/storage"
)

// Server is the object that listens to and handles all incoming HTTP requests.
type Server struct {
	cfg *Config
	mux *http.ServeMux
	db  storage.StatsKeeperStorage
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("ServeHTTP: recovered from panic: %v", r)
			writeErrorResponse(w, http.StatusInternalServerError, "server unavailable", fmt.Errorf("%v", r))
		}
	}()

	start := time.Now()
	rw := &responseWriter{
		status: http.StatusOK,
		actual: w,
	}
	s.mux.ServeHTTP(rw, r)
	end := time.Now()

	logrus.Infof("%s %s %s, %d %s, %d bytes", r.Method, r.URL.Path, end.Sub(start), rw.status, http.StatusText(rw.status), rw.bytesWritten)
}

// NewServer creates and initializes a new Server object using the given Config.
func NewServer(cfg *Config) (*Server, error) {
	db, err := storage.NewStatsKeeperStorage(cfg.DatabaseUrl)
	if err != nil {
		return nil, err
	}
	return &Server{
		cfg: cfg,
		db:  db,
	}, nil
}

// ListenHTTP initiates the HTTP listening and serving incoming requests. It returns only when process ends.
func (s *Server) ListenHTTP() error {
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/api/stats/list", s.ListUserStats)
	s.mux.HandleFunc("/api/stats/get", s.GetStat)
	s.mux.HandleFunc("/api/stats/add", s.AddStat)
	s.mux.HandleFunc("/api/stats/delete", s.DeleteStat)
	s.mux.HandleFunc("/api/stats/update", s.UpdateStat)

	logrus.Infof("serving http on :%d", s.cfg.HttpPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.HttpPort), s)
}
