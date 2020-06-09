package main

import (
	"net/http"
	"encoding/json"
)

type coaster struct {
	Name 			string `json:"name"`
	Manufacturer 	string `json:"manufacturer"`
	ID 				string `json:"id"`
	InPark 			string `json:"inPark"`
	Height 			int `json:"height"`
}

type coastersHandlers struct {
	store map[string]coaster
}

func newCoasterHandlers() *coastersHandlers {
	return &coastersHandlers{store: map[string]coaster{
		"id1": coaster{
			Name: "Fury 325",
			Manufacturer: "B+M",
			ID: "id1",
			InPark: "Luna Park",
			Height: 102,
		},
	}}
}

func (h *coastersHandlers) get(w http.ResponseWriter, r *http.Request) {
	coasters := make([]coaster, len(h.store))

	for i, coaster := range coasters {
		coasters[i] = coaster
	}

	json.NewEncoder(w).Encode(coasters)
}

func main() {

	coasterHandlers := newCoasterHandlers()

	http.HandleFunc("/coasters", coasterHandlers.get)
	http.ListenAndServe(":9000", nil)
}