package storage

// StatsKeeperStorage is the inteface that server will use to interact with the database.
type StatsKeeperStorage interface{}

// storage is the internal type that implements StatsKeeperStorage.
type storage struct{}

// NewStatsKeeperStorage creates a new StatsKeeperStorage by initializing connection to
// the database server at the given url. Any error during initialization is returned.
func NewStatsKeeperStorage(url string) (StatsKeeperStorage, error) {
	return &storage{}, nil
}
