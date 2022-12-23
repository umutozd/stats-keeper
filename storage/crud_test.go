package storage

import (
	"context"
	"os"
	"sort"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/umutozd/stats-keeper/protos/statspb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// newTestStorage creates a new storage instance, using the mongodb url in the environment
// or the default one if env url is empty. After creation, it drops collections and returns
// the storage instance.
func newTestStorage(t *testing.T) *storage {
	url := os.Getenv("TEST_MONGODB_URL")
	if url == "" {
		url = "mongodb://localhost:27017/?serverSelectionTimeoutMS=5000&connectTimeoutMS=10000"
	}

	s, err := NewStatsKeeperStorage(url)
	if err != nil {
		t.Fatalf("error creating new storage: %v", err)
	}

	ss := s.(*storage)
	if err = ss.statistics().Drop(context.Background()); err != nil {
		t.Fatalf("error dropping statistics collection: %v", err)
	}
	return ss
}

func Test_storage_CreateStatistic(t *testing.T) {
	tests := []struct {
		name   string
		entity *statspb.StatisticEntity
	}{
		{
			name: "case 1",
			entity: &statspb.StatisticEntity{
				Id:     "overriden-id-1",
				Name:   "entity-1",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Date{
					Date: &statspb.ComponentDate{
						Timestamps: []*timestamppb.Timestamp{
							{Seconds: 1, Nanos: 1},
							{Seconds: 2, Nanos: 2},
							{Seconds: 3, Nanos: 3},
						},
					},
				},
			},
		},
		{
			name: "case 2",
			entity: &statspb.StatisticEntity{
				Id:     "overriden-id-2",
				Name:   "entity-2",
				UserId: "user-2",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 123456,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStorage(t)
			result, err := s.CreateStatistic(context.TODO(), tt.entity)
			if err != nil {
				t.Fatalf("CreateStatistic returned unexpected error: %v", err)
			}

			// check if it's actually created
			got, err := s.GetStatistic(context.TODO(), result.Id)
			if err != nil {
				t.Fatalf("GetStatistic returned error after creating the entity: %v", err)
			}
			if diff := pretty.Compare(got, result); diff != "" {
				t.Fatalf("created and got are not the same, diff: %s", diff)
			}
		})
	}
}

func Test_storage_GetStatistic(t *testing.T) {

	tests := []struct {
		name           string
		entity         *statspb.StatisticEntity
		internalEntity *statisticEntity
		entityId       string
		expectedError  error
	}{
		{
			name:          "case 1 - not found",
			entity:        nil,
			entityId:      "entity-1",
			expectedError: NewErrorNotFound("statistic not found"),
		},
		{
			name: "case 2 - success",
			entity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 123,
					},
				},
			},
			entityId:      "id-1",
			expectedError: nil,
		},
		{
			name: "case 3 - extreme, empty id",
			entity: &statspb.StatisticEntity{
				Id:     "",
				Name:   "entity-2",
				UserId: "user-2",
				Component: &statspb.StatisticEntity_Date{
					Date: &statspb.ComponentDate{
						Timestamps: []*timestamppb.Timestamp{
							{Seconds: 123, Nanos: 456},
						},
					},
				},
			},
			entityId:      "",
			expectedError: nil,
		},
		{
			name:   "case 4 - deleted entity",
			entity: nil,
			internalEntity: &statisticEntity{
				Id:      "id-1",
				Name:    "entity-1",
				UserId:  "user-1",
				Counter: nil,
				Date: &statspb.ComponentDate{
					Timestamps: []*timestamppb.Timestamp{},
				},
				Deleted: true,
			},
			entityId:      "id-1",
			expectedError: NewErrorNotFound("statistic not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStorage(t)
			if tt.entity != nil {
				insertTestEntity(t, s, tt.entity)
			}
			if tt.internalEntity != nil {
				insertInternalTestEntity(t, s, tt.internalEntity)
			}

			got, err := s.GetStatistic(context.TODO(), tt.entityId)
			if hasError := compareErrors(t, tt.expectedError, err); hasError {
				return
			}
			if diff := pretty.Compare(got, tt.entity); diff != "" {
				t.Fatalf("wrong result, diff: %s", diff)
			}
		})
	}
}

func Test_storage_UpdateStatistic(t *testing.T) {
	tests := []struct {
		name   string
		entity *statspb.StatisticEntity

		fields         []string
		values         *statspb.StatisticEntity
		expectedError  error
		expectedEntity *statspb.StatisticEntity
	}{
		{
			name:   "case 1 - not found 1",
			entity: nil,
			fields: []string{"name", "counter"},
			values: &statspb.StatisticEntity{
				Id:   "id-1",
				Name: "name-1-updated",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 123,
					},
				},
			},
			expectedError:  NewErrorNotFound("statistic not found"),
			expectedEntity: nil,
		},
		{
			name:   "case 2 - not found 2",
			entity: nil,
			fields: []string{"id", "user_id"},
			values: &statspb.StatisticEntity{
				Id:     "id-1-updated",
				UserId: "user-1-updated",
			},
			expectedError:  NewErrorNotFound("statistic not found"),
			expectedEntity: nil,
		},
		{
			name: "case 3 - no update possible",
			entity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 123,
					},
				},
			},
			fields: []string{"invalid-field-1", "invalid-field-2", "invalid-field-3"},
			values: &statspb.StatisticEntity{
				Id: "id-1",
			},
			expectedError:  NewErrorNoUpdate("no update possible"),
			expectedEntity: nil,
		},
		{
			name: "case 4 - cannot change component type 1",
			entity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 123,
					},
				},
			},
			fields: []string{"date"},
			values: &statspb.StatisticEntity{
				Id: "id-1",
				Component: &statspb.StatisticEntity_Date{
					Date: &statspb.ComponentDate{
						Timestamps: []*timestamppb.Timestamp{},
					},
				},
			},
			expectedError:  NewErrorInvalidArgument("component cannot be changed from %s to %s", statspb.ComponentType_COUNTER, statspb.ComponentType_DATE),
			expectedEntity: nil,
		},
		{
			name: "case 5 - cannot change component type 2",
			entity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Date{
					Date: &statspb.ComponentDate{
						Timestamps: []*timestamppb.Timestamp{},
					},
				},
			},
			fields: []string{"counter"},
			values: &statspb.StatisticEntity{
				Id: "id-1",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 123,
					},
				},
			},
			expectedError:  NewErrorInvalidArgument("component cannot be changed from %s to %s", statspb.ComponentType_DATE, statspb.ComponentType_COUNTER),
			expectedEntity: nil,
		},
		{
			name: "case 6 - counter component update with date component",
			entity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 123,
					},
				},
			},
			fields: []string{"counter"},
			values: &statspb.StatisticEntity{
				Id: "id-1",
				Component: &statspb.StatisticEntity_Date{
					Date: &statspb.ComponentDate{
						Timestamps: []*timestamppb.Timestamp{},
					},
				},
			},
			expectedError:  NewErrorNoUpdate("no update possible"),
			expectedEntity: nil,
		},
		{
			name: "case 7 - date component update with counter component",
			entity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Date{
					Date: &statspb.ComponentDate{
						Timestamps: []*timestamppb.Timestamp{},
					},
				},
			},
			fields: []string{"date"},
			values: &statspb.StatisticEntity{
				Id: "id-1",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 123,
					},
				},
			},
			expectedError:  NewErrorNoUpdate("no update possible"),
			expectedEntity: nil,
		},
		{
			name: "case 8 - success 1",
			entity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Date{
					Date: &statspb.ComponentDate{
						Timestamps: []*timestamppb.Timestamp{},
					},
				},
			},
			fields: []string{"name", "date"},
			values: &statspb.StatisticEntity{
				Id:   "id-1",
				Name: "entity-1-updated",
				Component: &statspb.StatisticEntity_Date{
					Date: &statspb.ComponentDate{
						Timestamps: []*timestamppb.Timestamp{
							{Seconds: 1, Nanos: 1},
							{Seconds: 2, Nanos: 2},
							{Seconds: 3, Nanos: 3},
						},
					},
				},
			},
			expectedError: nil,
			expectedEntity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1-updated",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Date{
					Date: &statspb.ComponentDate{
						Timestamps: []*timestamppb.Timestamp{
							{Seconds: 1, Nanos: 1},
							{Seconds: 2, Nanos: 2},
							{Seconds: 3, Nanos: 3},
						},
					},
				},
			},
		},
		{
			name: "case 9 - success 2",
			entity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 3,
					},
				},
			},
			fields: []string{"name", "counter"},
			values: &statspb.StatisticEntity{
				Id:   "id-1",
				Name: "entity-1-updated",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 4,
					},
				},
			},
			expectedError: nil,
			expectedEntity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1-updated",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 4,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStorage(t)
			if tt.entity != nil {
				insertTestEntity(t, s, tt.entity)
			}

			got, err := s.UpdateStatistic(context.TODO(), tt.fields, tt.values)
			if hasError := compareErrors(t, tt.expectedError, err); hasError {
				return
			}
			if diff := pretty.Compare(got, tt.expectedEntity); diff != "" {
				t.Fatalf("wrong result, diff: %s", diff)
			}
		})
	}
}

func Test_storage_DeleteStatistic(t *testing.T) {
	tests := []struct {
		name     string
		entity   *statspb.StatisticEntity
		entityId string

		expectedError error
	}{
		{
			name:          "case 1 - not found",
			entity:        nil,
			entityId:      "id-1",
			expectedError: NewErrorNotFound("statistic not found"),
		},
		{
			name: "case 2 - success",
			entity: &statspb.StatisticEntity{
				Id:     "id-1",
				Name:   "entity-1",
				UserId: "user-1",
				Component: &statspb.StatisticEntity_Counter{
					Counter: &statspb.ComponentCounter{
						Count: 1,
					},
				},
			},
			entityId:      "id-1",
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStorage(t)
			if tt.entity != nil {
				insertTestEntity(t, s, tt.entity)
			}

			err := s.DeleteStatistic(context.TODO(), tt.entityId)
			if hasError := compareErrors(t, tt.expectedError, err); hasError {
				return
			}

			// check if it was actually deleted
			t.Log("checking if statistic is actually deleted")
			_, err = s.GetStatistic(context.TODO(), tt.entityId)
			_ = compareErrors(t, NewErrorNotFound("statistic not found"), err)
		})
	}
}

func Test_storage_ListUserStatistics(t *testing.T) {
	sortStatsSlice := func(slice []*statspb.StatisticEntity) []*statspb.StatisticEntity {
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].Id < slice[j].Id
		})
		return slice
	}

	tests := []struct {
		name     string
		entities []*statspb.StatisticEntity
		userId   string

		expectedResult []*statspb.StatisticEntity
	}{
		{
			name: "case 1 - no entities for user",
			entities: []*statspb.StatisticEntity{
				{
					Id:     "id-1",
					Name:   "entity-1",
					UserId: "user-1",
					Component: &statspb.StatisticEntity_Counter{
						Counter: &statspb.ComponentCounter{
							Count: 1,
						},
					},
				},
				{
					Id:     "id-2",
					Name:   "entity-2",
					UserId: "user-1",
					Component: &statspb.StatisticEntity_Counter{
						Counter: &statspb.ComponentCounter{
							Count: 1,
						},
					},
				},
				{
					Id:     "id-3",
					Name:   "entity-3",
					UserId: "user-1",
					Component: &statspb.StatisticEntity_Date{
						Date: &statspb.ComponentDate{
							Timestamps: []*timestamppb.Timestamp{
								{Seconds: 1, Nanos: 2},
							},
						},
					},
				},
			},
			userId:         "user-2",
			expectedResult: []*statspb.StatisticEntity{},
		},
		{
			name: "case 2 - some entities for user",
			entities: []*statspb.StatisticEntity{
				{
					Id:     "id-1",
					Name:   "entity-1",
					UserId: "user-1",
					Component: &statspb.StatisticEntity_Counter{
						Counter: &statspb.ComponentCounter{
							Count: 1,
						},
					},
				},
				{
					Id:     "id-2",
					Name:   "entity-2",
					UserId: "user-2",
					Component: &statspb.StatisticEntity_Counter{
						Counter: &statspb.ComponentCounter{
							Count: 1,
						},
					},
				},
				{
					Id:     "id-3",
					Name:   "entity-3",
					UserId: "user-2",
					Component: &statspb.StatisticEntity_Date{
						Date: &statspb.ComponentDate{
							Timestamps: []*timestamppb.Timestamp{
								{Seconds: 1, Nanos: 2},
							},
						},
					},
				},
			},
			userId: "user-2",
			expectedResult: sortStatsSlice([]*statspb.StatisticEntity{
				{
					Id:     "id-2",
					Name:   "entity-2",
					UserId: "user-2",
					Component: &statspb.StatisticEntity_Counter{
						Counter: &statspb.ComponentCounter{
							Count: 1,
						},
					},
				},
				{
					Id:     "id-3",
					Name:   "entity-3",
					UserId: "user-2",
					Component: &statspb.StatisticEntity_Date{
						Date: &statspb.ComponentDate{
							Timestamps: []*timestamppb.Timestamp{
								{Seconds: 1, Nanos: 2},
							},
						},
					},
				},
			}),
		},
		{
			name:           "case 3 - no entities at all",
			entities:       nil,
			userId:         "user-1",
			expectedResult: []*statspb.StatisticEntity{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStorage(t)
			for _, e := range tt.entities {
				insertTestEntity(t, s, e)
			}

			got, err := s.ListUserStatistics(context.TODO(), tt.userId)
			if hasError := compareErrors(t, nil, err); hasError {
				return
			}

			got = sortStatsSlice(got)
			if diff := pretty.Compare(got, tt.expectedResult); diff != "" {
				t.Fatalf("wrong result, diff: %s", diff)
			}
		})
	}
}

func insertTestEntity(t *testing.T, s *storage, entity *statspb.StatisticEntity) {
	se := &statisticEntity{}
	se.fromPB(entity)
	insertInternalTestEntity(t, s, se)
}

func insertInternalTestEntity(t *testing.T, s *storage, entity *statisticEntity) {
	if _, err := s.statistics().InsertOne(context.TODO(), entity); err != nil {
		t.Fatalf("error inserting test entity: %v", err)
	}
}

func compareErrors(t *testing.T, expected, got error) (hasError bool) {
	if got != nil {
		if expected == nil {
			t.Fatalf("expected nil error, but got non-nil error: %v", got)
		}
		if got.Error() != expected.Error() {
			t.Fatalf("wrong error: expected=%v, got=%v", expected, got)
		}
		return true
	} else {
		if expected != nil {
			t.Fatalf("got a nil error, but expected non-nil error: %v", expected)
		}
		return false
	}
}
