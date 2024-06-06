package statisticshandlers

import (
	"backend/internal/models"
	"backend/internal/services"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	StatisticsService services.StatisticsService
	JWTConfig         models.JWTConfig
}

func New(router fiber.Router, cfg Config, log *zerolog.Logger) {
	h := handler{
		service: cfg.StatisticsService,
		log:     log,
	}

	statisticsGroup := router.Group("/statistics", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(cfg.JWTConfig.JWTAccessSecretKey),
		},
	}))
	statisticsGroup.Get("/", h.get)
}
