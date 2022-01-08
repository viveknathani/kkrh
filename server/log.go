package server

import "net/http"

func (s *Server) handleLogStart(w http.ResponseWriter, r *http.Request)    {}
func (s *Server) handleLogEnd(w http.ResponseWriter, r *http.Request)      {}
func (s *Server) handleLogsPending(w http.ResponseWriter, r *http.Request) {}
