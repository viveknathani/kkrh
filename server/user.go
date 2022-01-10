package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/viveknathani/kkrh/entity"
	"github.com/viveknathani/kkrh/service"
	"github.com/viveknathani/kkrh/shared"
)

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {

	showRequestMetaData(s.Service.Logger, r)
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

	err = s.Service.Signup(r.Context(), &entity.User{
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

	showRequestMetaData(s.Service.Logger, r)
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

	token, err := s.Service.Login(r.Context(), &entity.User{
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
		HttpOnly: false,
		MaxAge:   int(time.Hour * 24 * 3),
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		Path:     "/",
		Secure:   true,
	})

	if ok := sendResponse(w, "ok", http.StatusOK); ok != nil {
		s.Service.Logger.Error(ok.Error())
	}
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {

	showRequestMetaData(s.Service.Logger, r)
	cookie, err := r.Cookie("token")
	if err != nil {
		s.Service.Logger.Error(err.Error())
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error())
		}
		return
	}

	err = s.Service.Logout(r.Context(), cookie.Value)
	if err != nil {

		s.Service.Logger.Error(err.Error())
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error())
		}
		return
	}

	if ok := sendResponse(w, "ok", http.StatusOK); ok != nil {
		s.Service.Logger.Error(ok.Error())
	}
}

func (s *Server) middlewareTokenVerification(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		showRequestMetaData(s.Service.Logger, r)
		cookie, err := r.Cookie("token")
		if err != nil {
			s.Service.Logger.Error(err.Error(), zapReqID(r))
			if ok := sendClientError(w, "not authenticated"); ok != nil {
				s.Service.Logger.Error(ok.Error(), zapReqID(r))
			}
			return
		}

		id, err := s.Service.VerifyAndDecodeToken(r.Context(), cookie.Value)
		if err != nil {

			s.Service.Logger.Error(err.Error(), zapReqID(r))
			if ok := sendClientError(w, "not authenticated"); ok != nil {
				s.Service.Logger.Error(ok.Error(), zapReqID(r))
			}
			return
		}

		request := r.Clone(shared.WithUserID(r.Context(), id))
		handler.ServeHTTP(w, request)
	}
}
