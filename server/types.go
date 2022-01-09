package server

// This file defines how the incoming and outgoing JSON payloads look like.

type genericResponse struct {
	Message string `json:"message"`
}

type userSignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type logStartRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Activity  string  `json:"activity"`
	StartTime int64   `json:"startTime"`
	Notes     string  `json:"notes"`
}

type logEndRequest struct {
	LogId   string `json:"logId"`
	EndTime int64  `json:"endTime"`
}
