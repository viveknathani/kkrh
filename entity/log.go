package entity

type Log struct {
	Id        string  `json:"id"`
	UserId    string  `json:"userId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Activity  string  `json:"activity"`
	StartTime int64   `json:"startTime"`
	EndTime   int64   `json:"endTime"`
	Notes     string  `json:"notes"`
}
