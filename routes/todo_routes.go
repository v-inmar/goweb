package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	handler_todo "github.com/v-inmar/goweb/handlers/todo_handlers"
)

func TodoRoutes(router *mux.Router, db *sql.DB, prefix string){
	subRouter := router.PathPrefix(prefix).Subrouter()

	// Get all todos
	subRouter.HandleFunc("", handler_todo.GetAllTodos(db)).Methods(http.MethodGet)
	// Create todo
	subRouter.HandleFunc("", handler_todo.CreateTodo(db)).Methods(http.MethodPost)
	// Get todo by id
	subRouter.HandleFunc("/{id}", handler_todo.GetTodo(db)).Methods(http.MethodGet)
	// Update todo by id
	subRouter.HandleFunc("/{id}", handler_todo.UpdateTodo(db)).Methods(http.MethodPut)
	// Delete a todo item by id
	subRouter.HandleFunc("/{id}", handler_todo.DeleteTodo(db)).Methods(http.MethodDelete)
}