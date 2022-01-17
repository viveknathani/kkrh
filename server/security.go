package server

import (
	"bytes"
	"net/http"
	"strconv"
)

type SecurityHandler struct {
	next http.Handler
}

func NewSecurityHandler(next http.Handler) *SecurityHandler {
	return &SecurityHandler{
		next: next,
	}
}

func createHSTSHeaderValue() string {
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

func (hsts *SecurityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("X-Frame-Options", "DENY")
	w.Header().Add("Content-Security-Policy", "default-src 'self'")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	if isHTTPS(r) {
		w.Header().Add("Strict-Transport-Security", createHSTSHeaderValue())
		hsts.next.ServeHTTP(w, r)
	} else {
		securedURL := "https://" + r.Host + r.RequestURI
		http.Redirect(w, r, securedURL, http.StatusMovedPermanently)
		return
	}
}
