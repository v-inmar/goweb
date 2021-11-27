package models

type ResponseBodyModel struct {
	PID   string `json:"pid"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}