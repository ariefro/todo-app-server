package routes

import (
	"github.com/ariefro/todo-app-server/handlers"
	"github.com/gofiber/fiber/v2"
)

func TodoRoutes(route fiber.Router) {
	api := route.Group("/api")

	todo := api.Group("/todo")

	todo.Get("/", handlers.GetTodos)
	todo.Post("/", handlers.CreateTodo)
	todo.Get("/:id", handlers.GetTodo)
	todo.Put("/:id", handlers.UpdateTodo)
}