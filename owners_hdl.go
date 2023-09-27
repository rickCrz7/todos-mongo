package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type OwnersHandler struct {
	dao OwnersDao
}

func NewOwnersHandler(dao OwnersDao, r *mux.Router) *OwnersHandler {
	h := &OwnersHandler{dao: dao}
	r.HandleFunc("/api/v1/owners", h.GetOwners).Methods("GET")
	r.HandleFunc("/api/v1/owners", h.CreateOwners).Methods("POST")
	r.HandleFunc("/api/v1/owners/{ownerID}", h.GetOwners).Methods("GET")
	r.HandleFunc("/api/v1/owners/{ownerID}", h.UpdateOwners).Methods("PUT")
	r.HandleFunc("/api/v1/owners/{ownerID}", h.DeleteOwners).Methods("DELETE")

	return h
}

func (h *OwnersHandler) GetOwners(w http.ResponseWriter, r *http.Request) {
	log.Println("GetOwners")
	owners, err := h.dao.GetAll()
	if err != nil {
		log.Println("Could not get owners: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(owners)
	w.WriteHeader(http.StatusOK)
}

func (h *OwnersHandler) CreateOwners(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateOwners")
	owner := &Owner{}
	err := json.NewDecoder(r.Body).Decode(owner)
	if err != nil {
		log.Println("Could not decode owner: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.dao.Create(owner)
	if err != nil {
		log.Println("Could not create owner: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *OwnersHandler) GetOwner(w http.ResponseWriter, r *http.Request) {
	log.Println("GetOwner")
	ownerID := mux.Vars(r)["ownerID"]
	owner, err := h.dao.Get(ownerID)
	if err != nil {
		log.Println("Could not get owner: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(owner)
	w.WriteHeader(http.StatusOK)
}

func (h *OwnersHandler) UpdateOwner(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateOwner")
	ownerID := mux.Vars(r)["ownerID"]
	owner := &Owner{}
	err := json.NewDecoder(r.Body).Decode(owner)
	if err != nil {
		log.Println("Could not decode owner: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	owner.ID = ownerID
	err = h.dao.Update(owner)
	if err != nil {
		log.Println("Could not update owner: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *OwnersHandler) DeleteOwner(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteOwner")
	ownerID := mux.Vars(r)["ownerID"]
	err := h.dao.Delete(ownerID)
	if err != nil {
		log.Println("Could not delete owner: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}