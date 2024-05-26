package usershandlers

import (
	"backend/internal/models"
	"backend/internal/services"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	UserService services.UserService
	JWTConfig   models.JWTConfig
}

func New(router fiber.Router, cfg Config, log *zerolog.Logger) {
	h := handler{
		service: cfg.UserService,
		log:     log,
	}

	userGroup := router.Group("/user", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(cfg.JWTConfig.JWTAccessSecretKey),
		},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			authorizationHeaderValue := ctx.Get(fiber.HeaderAuthorization)
			token := strings.Split(authorizationHeaderValue, "Bearer ")[1]
			if len(token) == 0 {
				return fiber.NewError(fiber.StatusUnauthorized, "missed jwt token")
			}

			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTConfig.JWTAccessSecretKey), nil
			})
			if err != nil {
				log.Error().Err(err).Send()
				return fiber.NewError(fiber.StatusUnauthorized, "err token parse")
			}

			ctx.Locals("claims", claims)
			return ctx.Next()
		},
	}))
	userGroup.Get("/", h.getList)
	userGroup.Get("/:id", h.getById)
}
