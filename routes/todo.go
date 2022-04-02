package routes

import (
	"github.com/ariefro/todo-app-server/handlers"
	"github.com/gofiber/fiber/v2"
)

func TodoRoutes(route fiber.Router) {
	todo := route.Group("/todo")

	todo.Get("/", handlers.GetTodos)
}