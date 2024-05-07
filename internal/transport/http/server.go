package http

import (
	"backend/internal/services"
	"backend/internal/transport/http/v1/answershandlers"
	"backend/internal/transport/http/v1/authhandlers"
	"backend/internal/transport/http/v1/fileshandlers"
	"backend/internal/transport/http/v1/groupshandlers"
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
	FileService   services.FileService
	GroupService  services.GroupService
	UserService   services.UserService
	Log           *zerolog.Logger
}

type Server struct {
	app  *fiber.App
	addr string

	taskService   services.TaskService
	answerService services.AnswerService
	fileService   services.FileService
	groupService  services.GroupService
	userService   services.UserService

	log *zerolog.Logger
}

func NewServer(cfg *Config) *Server {
	s := &Server{
		app:           nil,
		addr:          cfg.Addr,
		taskService:   cfg.TaskService,
		answerService: cfg.AnswerService,
		fileService:   cfg.FileService,
		groupService:  cfg.GroupService,
		userService:   cfg.UserService,
		log:           cfg.Log,
	}

	s.app = fiber.New(fiber.Config{
		ErrorHandler:          s.errorHandler,
		DisableStartupMessage: true,
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
	fileshandlers.New(v1Group, fileshandlers.Config{FileService: s.fileService}, s.log)
	groupshandlers.New(v1Group, groupshandlers.Config{GroupService: s.groupService}, s.log)
	authhandlers.New(v1Group, authhandlers.Config{UserService: s.userService}, s.log)
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
			if err = ctx.SendString(err.Error()); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong")
			}
			return nil
		}
	}

	return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong")
}
