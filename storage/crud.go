package storage

import (
	"context"

	"github.com/umutozd/stats-keeper/protos/statspb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *storage) CreateStatistic(ctx context.Context, entity *statspb.StatisticEntity) (*statspb.StatisticEntity, error) {
	se := &statisticEntity{}
	se.fromPB(entity)
	se.Id = primitive.NewObjectID().Hex()

	if _, err := s.statistics().InsertOne(ctx, se); err != nil {
		return nil, err
	}
	return se.toPB(), nil
}

func (s *storage) GetStatistic(ctx context.Context, entityId string) (*statspb.StatisticEntity, error) {
	se := &statisticEntity{}
	filter := bson.M{"_id": entityId}
	if err := s.statistics().FindOne(ctx, filter).Decode(se); err != nil {
		return nil, err
	}
	return se.toPB(), nil
}

func (s *storage) UpdateStatistic(ctx context.Context, fields []string, values *statspb.StatisticEntity) (*statspb.StatisticEntity, error) {
	se := &statisticEntity{}
	filter := bson.M{"_id": values.Id}
	set := bson.M{}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	for _, f := range fields {
		switch f {
		case "id", "user_id":
			// immutable fields, skip them
		case "name":
			set[f] = values.Name
		case "counter":
			if comp := values.GetCounter(); comp != nil {
				set[f] = comp
			}
		case "date":
			if comp := values.GetDate(); comp != nil {
				set[f] = comp
			}
		}
	}
	if len(set) == 0 {
		// no need to make ineffectual update, short-circuit here
		return nil, ErrNoUpdatePossible
	}
	update := bson.M{"$set": set}
	if err := s.statistics().FindOneAndUpdate(ctx, filter, update, opts).Decode(&se); err != nil {
		return nil, err
	}
	return se.toPB(), nil
}

func (s *storage) DeleteStatistic(ctx context.Context, entityId string) error {
	filter := bson.M{"_id": entityId}
	update := bson.M{
		"$set": bson.M{
			"deleted": true,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	return s.statistics().FindOneAndUpdate(ctx, filter, update, opts).Err()
}

func (s *storage) ListUserStatistics(ctx context.Context, userId string) ([]*statspb.StatisticEntity, error) {
	var result []*statspb.StatisticEntity
	var internalResult []*statisticEntity

	filter := bson.M{"user_id": userId}
	cursor, err := s.statistics().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &internalResult); err != nil {
		return nil, err
	}

	for _, se := range internalResult {
		result = append(result, se.toPB())
	}
	return result, nil
}
