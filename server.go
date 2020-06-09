package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gofiber/fiber"
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

func (h *coastersHandler) getCoasters(c *fiber.Ctx) {
	c.JSON(h.store)
}

func (h *coastersHandler) getCoaster(c *fiber.Ctx) {

	id, err := strconv.ParseInt(c.Params("id"), 10, 0)

	if err != nil {
		c.Status(http.StatusNotFound).JSON(error{true, "Invalid coaster id"})
		return
	}

	if coaster, ok := h.store[int(id)]; ok {
		c.JSON(coaster)
		return
	}
	c.Status(http.StatusNotFound).JSON(error{true, "Coaster ID not found"})
}

func (h *coastersHandler) addCoaster(c *fiber.Ctx) {

	var coaster coaster

	err := c.BodyParser(&coaster)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(error{true, "Invalid coaster data"})
		return
	}

	h.Lock()
	h.store[coaster.ID] = coaster
	defer h.Unlock()
	c.JSON(h.store[coaster.ID])
}

func setupRoutes(app *fiber.App) {

	coasterHandlers := newCoasterHandlers()

	app.Get("/coasters", coasterHandlers.getCoasters)
	app.Post("/coasters", coasterHandlers.addCoaster)
	app.Get("/coasters/:id", coasterHandlers.getCoaster)

}

func main() {

	app := fiber.New()
	setupRoutes(app)
	app.Listen(9000)

}
