package main

import (
	"log"

	"github.com/ariefro/todo-app-server/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	err = app.Listen(":8000")
	if err != nil {
		panic(err)
	}
}