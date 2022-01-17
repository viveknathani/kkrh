package server

import "net/http"

func (s *Server) SetupRoutes() {

	s.Router.HandleFunc("/api/user/signup/", setContentTypeJSON(s.handleSignup)).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/user/login/", setContentTypeJSON(s.handleLogin)).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/user/logout/", setContentTypeJSON(s.handleLogout)).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/log/start/", setContentTypeJSON(s.middlewareTokenVerification(s.handleLogStart))).Methods(http.MethodPost)
	s.Router.HandleFunc("/api/log/stop/", setContentTypeJSON(s.middlewareTokenVerification(s.handleLogEnd))).Methods(http.MethodPut)
	s.Router.HandleFunc("/api/log/pending", setContentTypeJSON(s.middlewareTokenVerification(s.handleLogsPending))).Methods(http.MethodGet)
	s.Router.HandleFunc("/health", setContentTypeJSON(s.showThatIAmAlive))
	s.Router.Use(setContentTypeFileFormat)
	s.setupWeb("web")
	s.Router.HandleFunc("/", s.serveIndex)
}
