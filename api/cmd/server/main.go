package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AyitiDev/Bornage/api/internal/config"
	"github.com/AyitiDev/Bornage/api/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1 Load strongly typed configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Fatal: failed to load configuration: %v", err)
	}

	// 2 Initialize PostgreSQL / PostGIS Connection Pool (Singleton Pattern)
	dbConn, err := db.GetInstance(cfg.Database)
	if err != nil {
		log.Printf("[WARN] Database connection warning: %v (running in disconnected mode)", err)
	} else {
		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}
		}()
	}

	// 3 Initialize Fiber app with custom configuration
	app := fiber.New(fiber.Config{
		AppName:      "Open Land Registry API v1.0",
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// 4 Register global middleware (Chain of Responsibility)
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS, PATCH",
	}))

	startTime := time.Now()

	// 5 Health Check endpoint with DB status (Task 17 & Task 19)
	app.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "disconnected"
		if dbConn != nil {
			ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
			defer cancel()
			if err := dbConn.PingContext(ctx); err == nil {
				dbStatus = "connected"
			} else {
				dbStatus = "unhealthy"
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      "ok",
			"system":      "open-land-registry",
			"environment": cfg.Server.Env,
			"database":    dbStatus,
			"uptime":      time.Since(startTime).String(),
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	// 6 Graceful shutdown setup
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-shutdownChan
		log.Printf("Received termination signal (%v). Initiating graceful shutdown...", sig)

		// Allow in-flight requests up to 10 seconds to finish
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Printf("Error during graceful shutdown: %v", err)
		}
	}()

	// 7 Start listening
	addr := ":" + cfg.Server.Port
	log.Printf("Open Land Registry API server running on port %s (env: %s)...", cfg.Server.Port, cfg.Server.Env)

	if err := app.Listen(addr); err != nil {
		log.Printf("Server shut down: %v", err)
	}

	log.Println("Server exited cleanly.")
}
