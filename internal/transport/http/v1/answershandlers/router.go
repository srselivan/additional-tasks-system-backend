package answershandlers

import (
	"backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	AnswerService services.AnswerService
}

func New(router fiber.Router, cfg Config, log *zerolog.Logger) {
	h := handler{
		service: cfg.AnswerService,
		log:     log,
	}

	answerGroup := router.Group("/answer")
	answerGroup.Get("/:id", h.getById)
	answerGroup.Get("/", h.getList)
	answerGroup.Post("/", h.create)
	answerGroup.Put("/:id", h.update)
	answerGroup.Delete("/:id", h.delete)
}
