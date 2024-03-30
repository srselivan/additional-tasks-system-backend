package http

import (
	"backend/internal/services"
	"backend/internal/transport/http/v1/answershandlers"
	"backend/internal/transport/http/v1/taskshandlers"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
)

type Config struct {
	Addr          string
	TaskService   services.TaskService
	AnswerService services.AnswerService
	Log           *zerolog.Logger
}

type Server struct {
	app  *fiber.App
	addr string

	taskService   services.TaskService
	answerService services.AnswerService

	log *zerolog.Logger
}

func NewServer(cfg *Config) *Server {
	s := &Server{
		app:           nil,
		addr:          cfg.Addr,
		taskService:   cfg.TaskService,
		answerService: cfg.AnswerService,
		log:           cfg.Log,
	}

	s.app = fiber.New(fiber.Config{
		ErrorHandler: s.errorHandler,
	})

	s.init()

	return s
}

func (s *Server) Run() error {
	return s.app.Listen(s.addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

func (s *Server) init() {
	s.app.Use(cors.New())

	apiGroup := s.app.Group("/api")

	v1Group := apiGroup.Group("/v1")
	taskshandlers.New(v1Group, taskshandlers.Config{TaskService: s.taskService}, s.log)
	answershandlers.New(v1Group, answershandlers.Config{AnswerService: s.answerService}, s.log)
}

func (s *Server) errorHandler(ctx *fiber.Ctx, err error) error {
	s.log.Error().
		Err(err).
		Str("url", ctx.OriginalURL()).
		Send()

	s.log.Debug().
		Err(err).
		Str("url", ctx.OriginalURL()).
		Any("query", ctx.Queries()).
		Any("params", ctx.AllParams()).
		Str("body", string(ctx.Body())).
		Send()

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		switch fiberError.Code {
		case fiber.StatusInternalServerError:
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong")
		default:
			return err
		}
	}

	return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong")
}
