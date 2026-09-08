package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/rbravo/interseguro-challenge/api-go/internal/config"
	"github.com/rbravo/interseguro-challenge/api-go/internal/httpapi"
	"github.com/rbravo/interseguro-challenge/api-go/internal/statsclient"
)

func main() {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		AppName:      "Interseguro QR API (Go + Fiber)",
		ErrorHandler: errorHandler,
	})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	handler := httpapi.NewHandler(cfg, statsclient.New(cfg.StatsAPIURL, cfg.StatsTimeout))
	httpapi.Register(app, handler, cfg.JWTSecret)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("cerrando el servidor...")
		_ = app.Shutdown()
	}()

	log.Printf("qr-api-go escuchando en :%s (stats: %s)", cfg.Port, cfg.StatsAPIURL)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("error al iniciar el servidor: %v", err)
	}
}

// errorHandler unifica el formato de error de toda la API.
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error": fiber.Map{
			"status":  code,
			"message": err.Error(),
			"path":    c.Path(),
		},
	})
}
