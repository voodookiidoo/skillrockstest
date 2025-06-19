package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"skillrockstest/internal/app"
	"skillrockstest/pkg/db"
)

func main() {
	f := fiber.New()
	pgConn := db.MustConnect()
	redConn := db.MustConnectRedis()
	a := app.NewApp(pgConn, redConn)

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
		Get("/tasks", a.GetAll).
		Get("/tasks/:id", a.Get).
		Post("/tasks", a.Post).
		Put("/tasks/:id", a.Put).
		Delete("/tasks/:id", a.Delete)

	if err := f.Listen(":8080"); err != nil {
		panic(err)
	}
}
