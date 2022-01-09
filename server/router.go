package server

import "net/http"

func (s *Server) SetupRoutes() {

	s.Router.HandleFunc("/api/user/signup/", s.handleSignup).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/user/login/", s.handleLogin).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/user/logout/", s.handleLogout).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/log/start/", s.middlewareTokenVerification(s.handleLogStart)).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/log/stop/", s.middlewareTokenVerification(s.handleLogEnd)).Methods(http.MethodPut)
	s.Router.HandleFunc("/api/log/pending", s.middlewareTokenVerification(s.handleLogsPending)).Methods(http.MethodGet)
}
