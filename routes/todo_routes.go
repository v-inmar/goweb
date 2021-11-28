package routes_todo

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	todos "github.com/v-inmar/goweb/handlers/todo_handlers"
)

func TodoRoutes(router *mux.Router, db *sql.DB, prefix string){
	subRouter := router.PathPrefix(prefix).Subrouter()

	// Get all todos
	subRouter.HandleFunc("", todos.GetAllTodos(db)).Methods(http.MethodGet)
	// Create todo
	subRouter.HandleFunc("", todos.CreateTodo(db)).Methods(http.MethodPost)
	// Get todo by id
	subRouter.HandleFunc("/{id}", todos.GetTodo(db)).Methods(http.MethodGet)
	// Update todo by id
	subRouter.HandleFunc("/{id}", todos.UpdateTodo(db)).Methods(http.MethodPut)
	// Delete a todo item by id
	subRouter.HandleFunc("/{id}", todos.DeleteTodo(db)).Methods(http.MethodDelete)
}