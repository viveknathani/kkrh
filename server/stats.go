package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/viveknathani/kkrh/shared"
)

func (s *Server) handleLogsInRange(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	startTime, err := strconv.ParseInt(params["startTime"][0], 10, 64)
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(err.Error(), zapReqID(r))
		}
		return
	}

	endTime, err := strconv.ParseInt(params["endTime"][0], 10, 64)
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(err.Error(), zapReqID(r))
		}
		return
	}

	list, err := s.Service.GetLogsInRange(r.Context(), shared.ExtractUserID(r.Context()), startTime, endTime)
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
	w.WriteHeader(http.StatusOK)
	if _, ok := w.Write(data); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
	showRequestEnd(s.Service.Logger, r)
}
