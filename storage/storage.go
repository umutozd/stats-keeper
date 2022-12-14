package storage

import (
	"context"

	"github.com/umutozd/stats-keeper/protos/statspb"
)

// StatsKeeperStorage is the inteface that server will use to interact with the database.
type StatsKeeperStorage interface {
	CreateStatistic(ctx context.Context, entity *statspb.StatisticEntity) (*statspb.StatisticEntity, error)
	GetStatistic(ctx context.Context, entityId string) (*statspb.StatisticEntity, error)
	UpdateStatistic(ctx context.Context, fields []string, values *statspb.StatisticEntity) (*statspb.StatisticEntity, error)
	DeleteStatistic(ctx context.Context, entityId string) error
	ListUserStatistics(ctx context.Context, userId string) ([]*statspb.StatisticEntity, error)
}

// storage is the internal type that implements StatsKeeperStorage.
type storage struct{}

// NewStatsKeeperStorage creates a new StatsKeeperStorage by initializing connection to
// the database server at the given url. Any error during initialization is returned.
func NewStatsKeeperStorage(url string) (StatsKeeperStorage, error) {
	return &storage{}, nil
}
