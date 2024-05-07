package app

import (
	"backend/config"
	"backend/internal/services/answer"
	answerpgrepo "backend/internal/services/answer/repo/pg"
	"backend/internal/services/file"
	filepgrepo "backend/internal/services/file/repo/pg"
	"backend/internal/services/group"
	grouppgrepo "backend/internal/services/group/repo/pg"
	"backend/internal/services/task"
	taskpgrepo "backend/internal/services/task/repo/pg"
	"backend/internal/services/user"
	userpgrepo "backend/internal/services/user/repo/pg"
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

	answersRepo := answerpgrepo.NewAnswersRepo(pgConn)
	filesRepo := filepgrepo.NewFilesRepo(pgConn)
	groupsRepo := grouppgrepo.NewGroupsRepo(pgConn)
	tasksRepo := taskpgrepo.NewTasksRepo(pgConn)
	usersRepo := userpgrepo.NewUsersRepo(pgConn)

	answerService := answer.New(answersRepo, log)
	fileService := file.New(filesRepo, log)
	groupService := group.New(groupsRepo, log)
	taskService := task.New(tasksRepo, log)
	userService := user.New(usersRepo, log)

	server := http.NewServer(&http.Config{
		Addr:          cfg.HTTPServer.Addr,
		TaskService:   taskService,
		AnswerService: answerService,
		FileService:   fileService,
		GroupService:  groupService,
		UserService:   userService,
		Log:           log,
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
