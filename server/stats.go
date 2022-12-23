package server

import (
	"net/http"

	"github.com/umutozd/stats-keeper/protos/statspb"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) ListUserStats(w http.ResponseWriter, r *http.Request) {
	in := unmarshalRequestBody(w, r, &statspb.ListUserStatisticsRequest{})
	if in == nil {
		return
	}
	if in.UserId == "" {
		writeErrorResponse(w, http.StatusBadRequest, "user_id cannot be empty", nil)
		return
	}

	entities, err := s.db.ListUserStatistics(r.Context(), in.UserId)
	if err != nil {
		writeStorageError(w, err)
		return
	}

	writeJsonResponse(w, http.StatusOK, &statspb.ListUserStatisticsResponse{Entities: entities})
}

func (s *Server) GetStat(w http.ResponseWriter, r *http.Request) {
	in := unmarshalRequestBody(w, r, &statspb.GetStatisticRequest{})
	if in == nil {
		return
	}
	if in.EntityId == "" {
		writeErrorResponse(w, http.StatusBadRequest, "entity_id cannot be empty", nil)
		return
	}

	entity, err := s.db.GetStatistic(r.Context(), in.EntityId)
	if err != nil {
		writeStorageError(w, err)
		return
	}

	writeJsonResponse(w, http.StatusOK, entity)
}

func (s *Server) AddStat(w http.ResponseWriter, r *http.Request) {
	in := unmarshalRequestBody(w, r, &statspb.StatisticEntity{})
	if in == nil {
		return
	}
	if in.Name == "" || in.UserId == "" || in.Component == nil {
		writeErrorResponse(w, http.StatusBadRequest, "name, user_id and component cannot be empty", nil)
		return
	}

	entity, err := s.db.CreateStatistic(r.Context(), in)
	if err != nil {
		writeStorageError(w, err)
		return
	}
	writeJsonResponse(w, http.StatusOK, entity)
}

func (s *Server) DeleteStat(w http.ResponseWriter, r *http.Request) {
	in := unmarshalRequestBody(w, r, &statspb.DeleteStatisticRequest{})
	if in == nil {
		return
	}
	if in.EntityId == "" {
		writeErrorResponse(w, http.StatusBadRequest, "entity_id cannot be empty", nil)
		return
	}

	if err := s.db.DeleteStatistic(r.Context(), in.EntityId); err != nil {
		writeStorageError(w, err)
		return
	} else {
		writeJsonResponse(w, http.StatusOK, &emptypb.Empty{})
	}
}

func (s *Server) UpdateStat(w http.ResponseWriter, r *http.Request) {
	in := unmarshalRequestBody(w, r, &statspb.UpdateStatisticRequest{})
	if in == nil {
		return
	}
	if in.Fields == nil || len(in.Fields.Paths) == 0 || in.Values == nil {
		writeErrorResponse(w, http.StatusBadRequest, "fields.paths and values must be non-empty or non-null", nil)
		return
	}

	entity, err := s.db.UpdateStatistic(r.Context(), in.Fields.Paths, in.Values)
	if err != nil {
		writeStorageError(w, err)
		return
	}
	writeJsonResponse(w, http.StatusOK, entity)
}
