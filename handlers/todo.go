package handlers

import (
	"fmt"
	"os"

	"github.com/ariefro/todo-app-server/config"
	"github.com/ariefro/todo-app-server/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTodos(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	query := bson.D{{}}

	cursor, err := todoCollection.Find(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": "Something went wrong",
			"error": err.Error(),
		})
	}

	var todos []models.Todo = make([]models.Todo, 0)

	fmt.Println(todos)
	fmt.Println(&todos)
	err = cursor.All(c.Context(), &todos)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"todos": todos,
		},
	})
}
