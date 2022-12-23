package storage

import (
	"fmt"
	"net/http"

	"github.com/umutozd/stats-keeper/protos/statspb"
)

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

type storageError struct {
	Message string
	Type    storageErrorType
}

type storageErrorType int8

const (
	storageErrorType_INVALID_ARGUMENT = 1
	storageErrorType_NOT_FOUND        = 2
	storageErrorType_NO_UPDATE        = 3
	storageErrorType_INTERNAL         = 4
)

func (set storageErrorType) String() string {
	switch set {
	case storageErrorType_INVALID_ARGUMENT:
		return "INVALID_ARGUMENT"
	case storageErrorType_NOT_FOUND:
		return "NOT_FOUND"
	case storageErrorType_NO_UPDATE:
		return "NO_UPDATE"
	case storageErrorType_INTERNAL:
		return "INTERNAL"
	default:
		return "UNKNOWN"
	}
}

func (set storageErrorType) HttpStatus() int {
	switch set {
	case storageErrorType_INVALID_ARGUMENT:
		return http.StatusBadRequest
	case storageErrorType_NOT_FOUND:
		return http.StatusNotFound
	case storageErrorType_NO_UPDATE:
		return http.StatusBadRequest
	case storageErrorType_INTERNAL:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func (se *storageError) Error() string {
	return fmt.Sprintf("storage: type=%s; %s", se.Type, se.Message)
}

func NewErrorInvalidArgument(format string, args ...any) error {
	return &storageError{Message: fmt.Sprintf(format, args...), Type: storageErrorType_INVALID_ARGUMENT}
}

func NewErrorNotFound(format string, args ...any) error {
	return &storageError{Message: fmt.Sprintf(format, args...), Type: storageErrorType_NOT_FOUND}
}

func NewErrorNoUpdate(format string, args ...any) error {
	return &storageError{Message: fmt.Sprintf(format, args...), Type: storageErrorType_NO_UPDATE}
}

func NewErrorInternal(format string, args ...any) error {
	return &storageError{Message: fmt.Sprintf(format, args...), Type: storageErrorType_INTERNAL}
}

func ToHttpError(err error) (code int, msg string) {
	if se, ok := err.(*storageError); ok {
		code, msg = se.Type.HttpStatus(), se.Message
	} else {
		code, msg = http.StatusInternalServerError, err.Error()
	}
	return
}
