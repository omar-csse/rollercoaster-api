package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type error struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type coaster struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	ID           string `json:"id"`
	InPark       string `json:"inPark"`
	Height       int    `json:"height"`
}

type coastersHandlers struct {
	sync.Mutex
	store map[string]coaster
}

func jsonERROR(w http.ResponseWriter, error interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(error)
}

func newCoasterHandlers() *coastersHandlers {
	return &coastersHandlers{store: map[string]coaster{
		"id1": {
			Name:         "Fury 325",
			Manufacturer: "B+M",
			ID:           "id1",
			InPark:       "Luna Park",
			Height:       102,
		},
	}}
}

func (h *coastersHandlers) coasters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		jsonERROR(w, error{true, "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
}

func (h *coastersHandlers) get(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.store)
}

func (h *coastersHandlers) post(w http.ResponseWriter, r *http.Request) {

	var coaster coaster

	err := json.NewDecoder(r.Body).Decode(&coaster)
	if err != nil {
		jsonERROR(w, error{true, "Invalid coaster body"}, http.StatusBadRequest)
		return
	}

	h.Lock()
	h.store[coaster.ID] = coaster
	defer h.Unlock()
}

func main() {

	coasterHandlers := newCoasterHandlers()

	http.HandleFunc("/coasters", coasterHandlers.coasters)
	http.ListenAndServe(":9000", nil)
}
