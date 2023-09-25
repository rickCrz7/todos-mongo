package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type TodoHandler struct {
	dao TodosDao
}

func NewTodoHandler(dao TodosDao, r *mux.Router) *TodoHandler {
	h := &TodoHandler{dao: dao}
	r.HandleFunc("/api/v1/todos", h.GetTodos).Methods("GET")
	r.HandleFunc("/api/v1/owners/{ownerID}/todos", h.CreateTodo).Methods("POST")
	r.HandleFunc("/api/v1/owners/{ownerID}/todos/{todoID}", h.GetTodo).Methods("GET")
	r.HandleFunc("/api/v1/owners/{ownerID}/todos/{todoID}", h.UpdateTodo).Methods("PUT")
	r.HandleFunc("/api/v1/owners/{ownerID}/todos/{todoID}", h.DeleteTodo).Methods("DELETE")

	return h
}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	log.Println("GetTodos")
	todos, err := h.dao.GetAll()
	if err != nil {
		log.Println("Could not get todos: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateTodo")
	ownerID := mux.Vars(r)["ownerID"]
	todo := &Todo{}
	err := json.NewDecoder(r.Body).Decode(todo)
	if err != nil {
		log.Println("Could not decode todo: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	todo.Owner = &Owner{ID: ownerID}
	err = h.dao.Create(todo)
	if err != nil {
		log.Println("Could not create todo: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	log.Println("GetTodo")
	todoID := mux.Vars(r)["todoID"]
	todo, err := h.dao.Get(todoID)
	if err != nil {
		log.Println("Could not get todo: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if todo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateTodo")
	todoID := mux.Vars(r)["todoID"]
	todo := &Todo{}
	err := json.NewDecoder(r.Body).Decode(todo)
	if err != nil {
		log.Println("Could not decode todo: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	todo.ID = todoID
	err = h.dao.Update(todo)
	if err != nil {
		log.Println("Could not update todo: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteTodo")
	todoID := mux.Vars(r)["todoID"]
	err := h.dao.Delete(todoID)
	if err != nil {
		log.Println("Could not delete todo: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
