package repo

import (
	"github.com/v-inmar/goweb/models"
)

var Todos = []models.TodoModel{
	{
		ID:    1,
		Title: "Wake up",
		Body:  "Open eyes and check the time",
	},
	{
		ID:    2,
		Title: "Make coffee",
		Body:  "Check the cupboard for clean mugs and use instant coffee",
	},
	{
		ID:    3,
		Title: "Go back to bed",
		Body:  "Check the time and go to sleep",
	},
}
