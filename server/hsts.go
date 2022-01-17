package server

import (
	"bytes"
	"net/http"
	"strconv"
)

type HandlerHSTS struct {
	next http.Handler
}

func NewHandlerHSTS(next http.Handler) *HandlerHSTS {
	return &HandlerHSTS{
		next: next,
	}
}

func createHeaderValue() string {
	buf := bytes.NewBufferString("max-age=")
	buf.WriteString(strconv.Itoa(63072000)) // 2y
	buf.WriteString("; includeSubDomains")
	return buf.String()
}

func isHTTPS(r *http.Request) bool {

	if r.Header.Get("X-Forwarded-Proto") == "https" {
		return true
	}

	if r.URL.Scheme == "https" {
		return true
	}

	if r.TLS != nil && r.TLS.HandshakeComplete {
		return true
	}

	return false
}

func (hsts *HandlerHSTS) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if isHTTPS(r) {
		w.Header().Add("Strict-Transport-Security", createHeaderValue())
		hsts.next.ServeHTTP(w, r)
	} else {
		securedURL := "https://" + r.Host + r.RequestURI
		http.Redirect(w, r, securedURL, http.StatusMovedPermanently)
		return
	}
}
