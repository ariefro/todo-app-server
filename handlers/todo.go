package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/ariefro/todo-app-server/config"
	"github.com/ariefro/todo-app-server/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
			"status": "error",
			"message": "cannot insert todo",
			"error": err.Error(),
		})
	}

	todo := &models.Todo{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}
	fmt.Println(bson.D{})

	todoCollection.FindOne(c.Context(), query).Decode(todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": todo,
		"data todo": fiber.Map{
			"todo": todo,
		},
	})
}

func GetTodo(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	paramId := c.Params("id")
	id, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"message": "cannont parse Id",
			"error": err.Error(),
		})
	}

	todo := &models.Todo{}
	
	query := bson.D{{Key: "_id", Value: id}}

	err = todoCollection.FindOne(c.Context(), query).Decode(todo)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"message": "todo not found",
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"todo": todo,
		},
	})
}

func UpdateTodo(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	paramId := c.Params("id")
	fmt.Println(paramId)

	id, err := primitive.ObjectIDFromHex(paramId)
	fmt.Println(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"message": "cannont parse Id",
			"error": err.Error(),
		})
	}

	data := new(models.Todo)
	err = c.BodyParser(data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"message": "cannot parse json",
			"error": err,
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

	var dataToUpdate bson.D

	if data.Title != nil {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "title", Value: data.Title})
	}

	if data.Completed != nil {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "completed", Value: data.Completed})
	}

	dataToUpdate = append(dataToUpdate, bson.E{Key: "updatedAt", Value: time.Now()})

	update := bson.D{{Key: "$set", Value: dataToUpdate}}

	err = todoCollection.FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"message": "todo not found",
				"error": err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"message": "cannot update todo",
			"error": err,
		})
	}

	todo := &models.Todo{}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"data": todo,
		},
	})
}

func DeleteTodo(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	paramId := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"message": "cannot parse json",
			"error": err,
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

	err = todoCollection.FindOneAndDelete(c.Context(), query).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"message": "todo not found",
				"error": err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"message": "cannot delete todo",
			"error": err,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}