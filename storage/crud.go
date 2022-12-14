package storage

import (
	"context"
	"errors"

	"github.com/umutozd/stats-keeper/protos/statspb"
)

func (s *storage) CreateStatistic(ctx context.Context, entity *statspb.StatisticEntity) (*statspb.StatisticEntity, error) {
	return nil, errors.New("not implemented")
}

func (s *storage) GetStatistic(ctx context.Context, entityId string) (*statspb.StatisticEntity, error) {
	return nil, errors.New("not implemented")
}

func (s *storage) UpdateStatistic(ctx context.Context, fields []string, values *statspb.StatisticEntity) (*statspb.StatisticEntity, error) {
	return nil, errors.New("not implemented")
}

func (s *storage) DeleteStatistic(ctx context.Context, entityId string) error {
	return errors.New("not implemented")
}

func (s *storage) ListUserStatistics(ctx context.Context, userId string) ([]*statspb.StatisticEntity, error) {
	return nil, errors.New("not implemented")
}
