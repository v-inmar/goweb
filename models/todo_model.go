package models

// No internal ids are shown to user
// only the public id, title and body
type PublicTodoModel struct {
	PID   string `json:"pid"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}
