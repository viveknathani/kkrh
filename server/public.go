package server

import (
	"net/http"
	"text/template"
)

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {

	p, err := template.ParseFiles("web/index.html")
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		return
	}
	err = p.Execute(w, nil)
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		return
	}
}
