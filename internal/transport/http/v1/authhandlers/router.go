package authhandlers

import (
	"backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	UserService services.UserService
}

func New(router fiber.Router, cfg Config, log *zerolog.Logger) {
	h := handler{
		service: cfg.UserService,
		log:     log,
	}

	router.Post("/sign-up", h.signUp)
}
