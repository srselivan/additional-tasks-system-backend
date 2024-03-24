package taskshandlers

import (
	"backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	TaskService services.TaskService
}

func New(router fiber.Router, cfg Config, log *zerolog.Logger) {
	h := handler{
		service: cfg.TaskService,
		log:     log,
	}

	taskGroup := router.Group("/task")
	taskGroup.Get("/:id", h.getById)
	taskGroup.Get("/", h.getList)
	taskGroup.Post("/", h.create)
	taskGroup.Put("/:id", h.update)
	taskGroup.Delete("/:id", h.delete)
}
