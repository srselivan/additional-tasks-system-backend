package groupshandlers

import (
	"backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	GroupService services.GroupService
}

func New(router fiber.Router, cfg Config, log *zerolog.Logger) {
	h := handler{
		service: cfg.GroupService,
		log:     log,
	}

	answerGroup := router.Group("/group")
	answerGroup.Get("/:id", h.getById)
	answerGroup.Get("/", h.getList)
	answerGroup.Post("/", h.create)
	answerGroup.Put("/:id", h.update)
	answerGroup.Delete("/:id", h.delete)
}
