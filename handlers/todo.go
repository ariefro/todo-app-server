package handlers

import (
	"fmt"
	"os"
	"time"

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

func CreateTodo(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	data := new(models.Todo)
	fmt.Println(data)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"message": "cannot parse JSON",
			"error": err.Error(),
		})
	}

	data.ID = nil
	f := false
	data.Completed = &f
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	result, err := todoCollection.InsertOne(c.Context(), data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "cannot insert todo",
			"error": err.Error(),
		})
	}

	todo := &models.Todo{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}
	fmt.Println(bson.D{})

	todoCollection.FindOne(c.Context(), query).Decode(todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data": todo,
		"data todo": fiber.Map{
			"todo": todo,
		},
	})
}