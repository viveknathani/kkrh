package server

func (s *Server) SetupRoutes() {

	s.Router.HandleFunc("/api/user/signup/", s.handleSignup).Methods("POST")
	s.Router.HandleFunc("/api/user/login/", s.handleLogin).Methods("POST")
	s.Router.HandleFunc("/api/user/logout/", s.handleLogout).Methods("POST")
	s.Router.HandleFunc("/api/log/start/", s.middlewareTokenVerification(s.handleLogStart)).Methods("POST")
	s.Router.HandleFunc("/api/log/stop/", s.middlewareTokenVerification(s.handleLogEnd)).Methods("PUT")
	s.Router.HandleFunc("/api/log/pending", s.middlewareTokenVerification(s.handleLogsPending)).Methods("GET")
}
