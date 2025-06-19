package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"skillrockstest/internal/app"
)

func main() {
	f := fiber.New()
	a := app.NewApp()
	defer a.Close()
	f.Use(logger.New(logger.ConfigDefault))

	f.
		Get("/tasks", a.GetAll).
		Get("/tasks/:id", a.Get).
		Post("/tasks", a.Post).
		Put("/tasks/:id", a.Put).
		Delete("/tasks/:id", a.Delete)

	if err := f.Listen(":8080"); err != nil {
		panic(err)
	}
}
