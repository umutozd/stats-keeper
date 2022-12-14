package server

import (
	"fmt"
	"net/http"

	"github.com/umutozd/stats-keeper/storage"
)

// Server is the object that listens to and handles all incoming HTTP requests.
type Server struct {
	cfg *Config
	db  storage.StatsKeeperStorage
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
	mux := http.NewServeMux()
	mux.HandleFunc("/api/stats/list", s.ListUserStats)
	mux.HandleFunc("/api/stats/get", s.GetStat)
	mux.HandleFunc("/api/stats/add", s.AddStat)
	mux.HandleFunc("/api/stats/delete", s.DeleteStat)
	mux.HandleFunc("/api/stats/update", s.UpdateStat)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.HttpPort), mux)
}
