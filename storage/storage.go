package storage

import (
	"context"

	"github.com/umutozd/stats-keeper/protos/statspb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultDatabaseName      = "StatsKeeper"
	statisticsCollectionName = "Statistics"
)

// StatsKeeperStorage is the inteface that server will use to interact with the database.
type StatsKeeperStorage interface {
	// CreateStatistic inserts the given entity to the database after initializing some of its data, such as Id.
	CreateStatistic(ctx context.Context, entity *statspb.StatisticEntity) (*statspb.StatisticEntity, error)

	// GetStatistic finds and returns the entity specified by entityId.
	GetStatistic(ctx context.Context, entityId string) (*statspb.StatisticEntity, error)

	// UpdateStatistic updates the entity specified by values.Id, using fields. Each element in fields specify which
	// field to update in the entity. Immutable fields such as Id or UserId are ignored. If no possible update is found,
	// ErrNoUpdatePossible is returned.
	UpdateStatistic(ctx context.Context, fields []string, values *statspb.StatisticEntity) (*statspb.StatisticEntity, error)

	// DeleteStatistic deletes the entity from database so that it cannot be found by any other CRUD method.
	DeleteStatistic(ctx context.Context, entityId string) error

	// ListUserStatistics returns a slice of entities belonging to the user specified by userId.
	ListUserStatistics(ctx context.Context, userId string) ([]*statspb.StatisticEntity, error)
}

// storage is the internal type that implements StatsKeeperStorage.
type storage struct {
	cli *mongo.Client
}

// NewStatsKeeperStorage creates a new StatsKeeperStorage by initializing connection to
// the database server at the given url. Any error during initialization is returned.
func NewStatsKeeperStorage(url string) (StatsKeeperStorage, error) {
	cli, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	return &storage{
		cli: cli,
	}, nil
}

// statistics returns a handle to the statistics collection in MongoDB
func (s *storage) statistics() *mongo.Collection {
	return s.cli.Database(defaultDatabaseName).Collection(statisticsCollectionName)
}
