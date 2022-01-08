package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/viveknathani/kkrh/service"
)

// Server holds together all the configuration needed to run this web service.
type Server struct {
	Service *service.Service
	Router  *mux.Router
}

// RequestID will be used in context
type RequestID string

// UserID will be used in context
type UserID string

// ServeHTTP is implemented so that Server can be used for listening to requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestID := uuid.New().String()
	request := r.Clone(context.WithValue(context.Background(), RequestID("requestID"), requestID))
	s.Router.ServeHTTP(w, request)
}
