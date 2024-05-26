package app

import (
	"backend/config"
	"backend/internal/models"
	repos "backend/internal/repo/pg"
	"backend/internal/services"
	"backend/internal/transport/http"
	"backend/pkg/logger"
	"backend/pkg/postgres"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Logger.Level)

	pgConn, err := postgres.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Open postgres connection error")
	}
	//if err = postgres.RunMigrations(pgConn.DB); err != nil {
	//	log.Fatal().Err(err).Msg("Run postgres migrations error")
	//}

	if _, err = os.Stat("storage"); os.IsNotExist(err) {
		if err = os.Mkdir("storage", 0777); err != nil {
			panic(err)
		}
	}

	answersRepo := repos.NewAnswersRepo(pgConn)
	filesRepo := repos.NewFilesRepo(pgConn)
	groupsRepo := repos.NewGroupsRepo(pgConn)
	tasksRepo := repos.NewTasksRepo(pgConn)
	taskLinksRepo := repos.NewTaskLinksRepo(pgConn)
	usersRepo := repos.NewUsersRepo(pgConn)
	authRepo := repos.NewAuthRepo(pgConn)
	marksRepo := repos.NewMarksRepo(pgConn)

	fileService := services.NewFileServiceImpl(filesRepo, log)
	answerService := services.NewAnswerServiceImpl(answersRepo, fileService, log)
	groupService := services.NewGroupServiceImpl(groupsRepo, log)
	taskLinksService := services.NewTaskLinksServiceImpl(taskLinksRepo, log)
	taskService := services.NewTaskServiceImpl(tasksRepo, fileService, taskLinksService, log)
	userService := services.NewUserServiceImpl(usersRepo, log)
	marksService := services.NewMarkServiceImpl(marksRepo, log)
	authService := services.NewAuthServiceImpl(
		authRepo,
		models.JWTConfig{
			JWTAccessExpirationTime:  cfg.JWT.JWTAccessTokenExpTime,
			JWTRefreshExpirationTime: cfg.JWT.JWTRefreshTokenExpTime,
			JWTAccessSecretKey:       cfg.JWT.JWTAccessSecretKey,
			JWTRefreshSecretKey:      cfg.JWT.JWTRefreshSecretKey,
		},
		userService,
		log,
	)

	server := http.NewServer(&http.Config{
		Addr:          cfg.HTTPServer.Addr,
		TaskService:   taskService,
		AnswerService: answerService,
		FileService:   fileService,
		GroupService:  groupService,
		UserService:   userService,
		AuthService:   authService,
		MarkService:   marksService,
		JWTConfig: models.JWTConfig{
			JWTAccessExpirationTime:  cfg.JWT.JWTAccessTokenExpTime,
			JWTRefreshExpirationTime: cfg.JWT.JWTRefreshTokenExpTime,
			JWTAccessSecretKey:       cfg.JWT.JWTAccessSecretKey,
			JWTRefreshSecretKey:      cfg.JWT.JWTRefreshSecretKey,
		},
		Log: log,
	})

	go func() {
		if err = server.Run(); err != nil {
			log.Fatal().Err(err).Msg("Start http server error")
		}
	}()
	log.Info().Msgf("Successfully run http server on %s", cfg.HTTPServer.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGKILL)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Shutdown http server error")
	}
	log.Info().Msg("Http server successfully stopped")

	if err = pgConn.Close(); err != nil {
		log.Fatal().Err(err).Msg("Close postgres connection error")
	}
	log.Info().Msg("Postgres connection successfully closed")
}
