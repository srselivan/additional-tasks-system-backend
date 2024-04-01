package answershandlers

import (
	"backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	FileService services.FileService
}

func New(router fiber.Router, cfg Config, log *zerolog.Logger) {
	h := handler{
		service: cfg.FileService,
		log:     log,
	}

	answerGroup := router.Group("/file")
	answerGroup.Get("/:id", h.getById)
	answerGroup.Post("/", h.create)
	answerGroup.Delete("/:id", h.delete)
}
