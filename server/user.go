package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/viveknathani/kkrh/entity"
	"github.com/viveknathani/kkrh/service"
)

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var u userSignupRequest
	err := decoder.Decode(&u)

	if err != nil {
		s.Service.Logger.Error(err.Error())
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error())
		}
		return
	}

	err = s.Service.Signup(&entity.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: []byte(u.Password),
	})
	if err != nil {

		s.Service.Logger.Error(err.Error())
		switch {
		case err == service.ErrEmailExists || err == service.ErrInvalidEmailFormat || err == service.ErrInvalidPasswordFormat:
			{
				if ok := sendClientError(w, err.Error()); ok != nil {
					s.Service.Logger.Error(ok.Error())
				}
				return
			}
		default:
			{
				if ok := sendServerError(w); ok != nil {
					s.Service.Logger.Error(ok.Error())
				}
				return
			}
		}
	}

	if ok := sendCreated(w); ok != nil {
		s.Service.Logger.Error(ok.Error())
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var u userLoginRequest
	err := decoder.Decode(&u)

	if err != nil {
		s.Service.Logger.Error(err.Error())
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error())
		}
		return
	}

	token, err := s.Service.Login(&entity.User{
		Email:    u.Email,
		Password: []byte(u.Password),
	})
	if err != nil {

		s.Service.Logger.Error(err.Error())
		switch {
		case err == service.ErrInvalidEmailPassword:
			{
				if ok := sendClientError(w, err.Error()); ok != nil {
					s.Service.Logger.Error(ok.Error())
				}
				return
			}
		default:
			{
				if ok := sendServerError(w); ok != nil {
					s.Service.Logger.Error(ok.Error())
				}
				return
			}
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})

	if ok := sendResponse(w, "ok", 200); ok != nil {
		s.Service.Logger.Error(ok.Error())
	}
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")
	if err != nil {
		s.Service.Logger.Error(err.Error())
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error())
		}
		return
	}

	err = s.Service.Logout(cookie.Value)
	if err != nil {

		s.Service.Logger.Error(err.Error())
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error())
		}
		return
	}

	if ok := sendResponse(w, "ok", 200); ok != nil {
		s.Service.Logger.Error(ok.Error())
	}
}

func (s *Server) middlewareTokenVerification(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			s.Service.Logger.Error(err.Error())
			if ok := sendServerError(w); ok != nil {
				s.Service.Logger.Error(ok.Error())
			}
			return
		}

		id, err := s.Service.VerifyAndDecodeToken(cookie.Value)
		if err != nil {

			s.Service.Logger.Error(err.Error())
			if ok := sendClientError(w, "not authenticated"); ok != nil {
				s.Service.Logger.Error(ok.Error())
			}
			return
		}

		request := r.Clone(context.WithValue(r.Context(), UserID("userId"), id))
		handler.ServeHTTP(w, request)
	}
}