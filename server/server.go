package server

import (
	"errors"

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
	return errors.New("not implemented")
}
