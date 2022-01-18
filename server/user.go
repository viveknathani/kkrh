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

	if true {
		if ok := sendClientError(w, "i wish you could go further"); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	decoder := json.NewDecoder(r.Body)
	var u userSignupRequest
	err := decoder.Decode(&u)

	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	err = s.Service.Signup(r.Context(), &entity.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: []byte(u.Password),
	})
	if err != nil {

		s.Service.Logger.Error(err.Error(), zapReqID(r))
		switch {
		case err == service.ErrEmailExists || err == service.ErrInvalidEmailFormat || err == service.ErrInvalidPasswordFormat:
			{
				if ok := sendClientError(w, err.Error()); ok != nil {
					s.Service.Logger.Error(ok.Error(), zapReqID(r))
				}
				return
			}
		default:
			{
				if ok := sendServerError(w); ok != nil {
					s.Service.Logger.Error(ok.Error(), zapReqID(r))
				}
				return
			}
		}
	}

	showRequestEnd(s.Service.Logger, r)
	if ok := sendCreated(w); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var u userLoginRequest
	err := decoder.Decode(&u)

	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	token, err := s.Service.Login(r.Context(), &entity.User{
		Email:    u.Email,
		Password: []byte(u.Password),
	})
	if err != nil {

		s.Service.Logger.Error(err.Error(), zapReqID(r))
		switch {
		case err == service.ErrInvalidEmailPassword:
			{
				if ok := sendClientError(w, err.Error()); ok != nil {
					s.Service.Logger.Error(ok.Error(), zapReqID(r))
				}
				return
			}
		default:
			{
				if ok := sendServerError(w); ok != nil {
					s.Service.Logger.Error(ok.Error(), zapReqID(r))
				}
				return
			}
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		MaxAge:   int(time.Hour.Seconds() * 24 * 3),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   (os.Getenv("MODE") == "prod"),
	})

	if ok := sendResponse(w, "ok", http.StatusOK); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
	showRequestEnd(s.Service.Logger, r)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	err = s.Service.Logout(r.Context(), cookie.Value)
	if err != nil {

		s.Service.Logger.Error(err.Error())
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	if ok := sendResponse(w, "ok", http.StatusOK); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
	showRequestEnd(s.Service.Logger, r)
}

func (s *Server) middlewareTokenVerification(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
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
