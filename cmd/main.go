package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"skillrockstest/internal/app"
	"skillrockstest/pkg/db"
)

func main() {
	f := fiber.New()
	conn := db.MustConnect()
	defer conn.Close(context.Background())
	a := app.NewApp(conn)
	
	//f.Use(requestLogger)
	// регаем мидлварь
	
	// регаем логгер
	f.Use(logger.New(logger.ConfigDefault))
	// регаем метрики
	requestsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total requests",
		Help: "Total number of requests",
	})
	prometheus.MustRegister(requestsTotal)
	promhttp.Handler()
	// регаем обработчики запросов
	
	f.
		Get("/tasks", a.Get).
		Post("/tasks", a.Post).
		Put("/tasks/:id", a.Put).
		Delete("/tasks/:id", a.Delete)
	
	if err := f.Listen(":8080"); err != nil {
		panic(err)
	}
}
