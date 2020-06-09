package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"sync"
)

type error struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type coaster struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	ID           int    `json:"id"`
	InPark       string `json:"inPark"`
	Height       int    `json:"height"`
}

type coastersHandler struct {
	sync.Mutex
	store map[int]coaster
}

func jsonERROR(w http.ResponseWriter, error interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(error)
}

func newCoasterHandlers() *coastersHandler {
	return &coastersHandler{store: map[int]coaster{
		1: {
			Name:         "Fury 325",
			Manufacturer: "B+M",
			ID:           1,
			InPark:       "Luna Park",
			Height:       102,
		},
	}}
}

func (h *coastersHandler) coasters(w http.ResponseWriter, r *http.Request) {
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

func (h *coastersHandler) get(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.store)
}

func (h *coastersHandler) post(w http.ResponseWriter, r *http.Request) {

	var coaster coaster

	err := json.NewDecoder(r.Body).Decode(&coaster)
	if err != nil {
		jsonERROR(w, error{true, "Invalid coaster body"}, http.StatusBadRequest)
		return
	}

	h.Lock()
	h.store[coaster.ID] = coaster
	defer h.Unlock()
	json.NewEncoder(w).Encode(h.store[coaster.ID])
}

func (h *coastersHandler) getCoaster(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(path.Base(r.URL.String()), 10, 0)

	if err != nil {
		jsonERROR(w, error{true, "Invalid coaster id"}, http.StatusBadRequest)
		return
	}

	if coaster, ok := h.store[int(id)]; ok {
		json.NewEncoder(w).Encode(coaster)
		return
	}
	jsonERROR(w, error{true, "coaster not found"}, http.StatusBadRequest)
}

func main() {

	coasterHandlers := newCoasterHandlers()

	http.HandleFunc("/coasters", coasterHandlers.coasters)
	http.HandleFunc("/coaster/", coasterHandlers.getCoaster)

	http.ListenAndServe(":9000", nil)
}
