package storage

import (
	"context"
	"errors"

	"github.com/umutozd/stats-keeper/protos/statspb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *storage) CreateStatistic(ctx context.Context, entity *statspb.StatisticEntity) (*statspb.StatisticEntity, error) {
	se := &statisticEntity{}
	se.fromPB(entity)
	se.Id = primitive.NewObjectID().Hex()

	if _, err := s.statistics().InsertOne(ctx, se); err != nil {
		return nil, NewErrorInternal(err, "error creating statistic")
	}
	return se.toPB(), nil
}

func (s *storage) GetStatistic(ctx context.Context, entityId string) (*statspb.StatisticEntity, error) {
	se := &statisticEntity{}
	filter := bson.M{"_id": entityId}
	if err := s.statistics().FindOne(ctx, filter).Decode(se); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, NewErrorNotFound(nil, "statistic not found")
		}
		return nil, NewErrorInternal(err, "error getting statistic from database")
	}
	if se.Deleted {
		return nil, NewErrorNotFound(nil, "statistic not found")
	}
	return se.toPB(), nil
}

func (s *storage) UpdateStatistic(ctx context.Context, fields []string, values *statspb.StatisticEntity) (*statspb.StatisticEntity, error) {
	se := &statisticEntity{}
	filter := bson.M{"_id": values.Id}
	set := bson.M{}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	entity, err := s.GetStatistic(ctx, values.Id)
	if err != nil {
		return nil, err
	}
	compType := entity.GetComponentType()

	for _, f := range fields {
		switch f {
		case "id", "user_id":
			return nil, NewErrorInvalidArgument(nil, "fields 'id', 'user_id' cannot be modified")
		case "name":
			set[f] = values.Name
		case "counter":
			if compType != statspb.ComponentType_COUNTER {
				return nil, NewErrorInvalidArgument(nil, "component cannot be changed from %s to %s", compType, statspb.ComponentType_COUNTER)
			}
			if comp := values.GetCounter(); comp != nil {
				set[f] = comp
			}
		case "date":
			if compType != statspb.ComponentType_DATE {
				return nil, NewErrorInvalidArgument(nil, "component cannot be changed from %s to %s", compType, statspb.ComponentType_DATE)
			}
			if comp := values.GetDate(); comp != nil {
				set[f] = comp
			}
		}
	}
	if len(set) == 0 {
		// no need to make ineffectual update, short-circuit here
		return nil, NewErrorNoUpdate(nil, "no update possible")
	}
	update := bson.M{"$set": set}
	if err := s.statistics().FindOneAndUpdate(ctx, filter, update, opts).Decode(&se); err != nil {
		return nil, NewErrorInternal(err, "error updating statistic")
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
	if err := s.statistics().FindOneAndUpdate(ctx, filter, update, opts).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NewErrorNotFound(nil, "statistic not found")
		}
		return NewErrorInternal(err, "error deleting statistic")
	}
	return nil
}

func (s *storage) ListUserStatistics(ctx context.Context, userId string) ([]*statspb.StatisticEntity, error) {
	var result []*statspb.StatisticEntity
	var internalResult []*statisticEntity

	filter := bson.M{"user_id": userId, "deleted": false}
	cursor, err := s.statistics().Find(ctx, filter)
	if err != nil {
		return nil, NewErrorInternal(err, "error listing statistics")
	}
	if err = cursor.All(ctx, &internalResult); err != nil {
		return nil, NewErrorInternal(err, "error decoding statistics")
	}

	for _, se := range internalResult {
		result = append(result, se.toPB())
	}
	if result == nil {
		result = []*statspb.StatisticEntity{}
	}
	return result, nil
}
