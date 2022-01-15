package server

import "net/http"

func (s *Server) showThatIAmAlive(w http.ResponseWriter, r *http.Request) {

	if ok := sendResponse(w, "alive", 200); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
}
