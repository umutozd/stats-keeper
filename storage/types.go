package storage

import "github.com/umutozd/stats-keeper/protos/statspb"

// statisticEntity is the internal representation of statspb.StatisticEntity. We need this type
// because statspb.StatisticEntity as a "oneof" field, which breaks MongoDB's marshal/unmarshal
// logic.
type statisticEntity struct {
	Id     string `bson:"_id"`
	Name   string `bson:"name"`
	UserId string `bson:"user_id"`

	Counter *statspb.ComponentCounter `bson:"counter"`
	Date    *statspb.ComponentDate    `bson:"date"`

	// Deleted reports whether this entity is deleted via a db call. Instead of actual delete, this
	// entity is marked as deleted. We may use this non-deleted entity in the future.
	Deleted bool `bson:"deleted"`
}

// toPB converts this statisticEntity to *statspb.StatisticEntity.
func (se *statisticEntity) toPB() *statspb.StatisticEntity {
	out := &statspb.StatisticEntity{
		Id:     se.Id,
		Name:   se.Name,
		UserId: se.UserId,
	}
	if se.Counter != nil {
		out.Component = &statspb.StatisticEntity_Counter{Counter: se.Counter}
	} else if se.Date != nil {
		out.Component = &statspb.StatisticEntity_Date{Date: se.Date}
	}
	return out
}

// fromPB converts the given *statspb.StatisticEntity to statisticEntity.
func (se *statisticEntity) fromPB(in *statspb.StatisticEntity) {
	se.Id = in.Id
	se.Name = in.Name
	se.UserId = in.UserId

	switch comp := in.Component.(type) {
	case *statspb.StatisticEntity_Counter:
		se.Counter = comp.Counter
	case *statspb.StatisticEntity_Date:
		se.Date = comp.Date
	}
}
