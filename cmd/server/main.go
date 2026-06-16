package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"users-api/config"
	db "users-api/db/sqlc"
	"users-api/internal/handler"
	"users-api/internal/logger"
	"users-api/internal/repository"
	"users-api/internal/routes"
	"users-api/internal/service"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger.Init(os.Getenv("APP_ENV") == "development")
	defer logger.Sync()

	sqlDB, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	logger.Log.Info("database connected")

	// Wire: SQLC → repository → service → handler
	queries  := db.New(sqlDB)
	userRepo := repository.NewUserRepo(queries)
	userSvc  := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	routes.Register(app, userHandler)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Log.Info("server starting", zap.String("addr", addr))

	if err := app.Listen(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
