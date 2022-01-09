package server

import (
	"encoding/json"
	"net/http"

	"github.com/viveknathani/kkrh/entity"
	"github.com/viveknathani/kkrh/service"
	"github.com/viveknathani/kkrh/shared"
)

func (s *Server) handleLogStart(w http.ResponseWriter, r *http.Request) {

	showRequestMetaData(s.Service.Logger, r)
	var l logStartRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&l)
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error())
		}
		return
	}

	err = s.Service.StartLog(r.Context(), &entity.Log{
		UserId:    shared.ExtractUserID(r.Context()),
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Activity:  l.Activity,
		StartTime: l.StartTime,
		Notes:     l.Notes,
	})

	if err != nil {
		switch {
		case err == service.ErrInvalidLog:
			{
				s.Service.Logger.Error(err.Error(), zapReqID(r))
				if ok := sendClientError(w, err.Error()); ok != nil {
					s.Service.Logger.Error(ok.Error())
				}
				return
			}
		default:
			{
				if ok := sendServerError(w); ok != nil {
					s.Service.Logger.Error(err.Error(), zapReqID(r))
				}
				return
			}
		}
	}

	if ok := sendCreated(w); ok != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
	}
}

func (s *Server) handleLogEnd(w http.ResponseWriter, r *http.Request) {

	showRequestMetaData(s.Service.Logger, r)
	var l logEndRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&l)
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error())
		}
		return
	}

	err = s.Service.EndLog(r.Context(), l.LogId, l.EndTime)

	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(err.Error(), zapReqID(r))
		}
		return
	}

	if ok := sendUpdated(w); ok != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
	}
}

func (s *Server) handleLogsPending(w http.ResponseWriter, r *http.Request) {

	showRequestMetaData(s.Service.Logger, r)
	list, err := s.Service.GetPendingLogs(r.Context(), shared.ExtractUserID(r.Context()))
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(err.Error(), zapReqID(r))
		}
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		return
	}
	w.WriteHeader(200)
	if _, ok := w.Write(data); ok != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
	}
}
