package models

type TodoModel struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// No internal ids are shown to user
// only the public id, title and body
type PublicTodoModel struct {
	PID   string `json:"pid"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
