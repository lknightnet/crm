package internal

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"project-service/config"
	"project-service/internal/controller"
	"project-service/internal/model"
	"project-service/internal/repository"
	"project-service/internal/service"
	"project-service/pkg/database"
	"project-service/pkg/server"
	"syscall"
)

func Run(cfg *config.Config) {
	connection, err := database.NewConnection(cfg.Database.URI, 0, 5, &model.Project{}, &model.ProjectUsers{},
		&model.ProjectUsers{}, &model.TaskExecutor{}, &model.Task{}, &model.TimerEntry{}, &model.InformationList{})
	if err != nil {
		panic(err)
	}

	repositories := repository.NewRepository(connection)

	deps := &service.DependenciesService{
		UserServiceApi:            cfg.UserService.Api,
		UserServicePort:           cfg.UserService.Port,
		ProjectRepository:         repositories.ProjectRepository,
		ProjectUsersRepository:    repositories.ProjectUsersRepository,
		TaskRepository:            repositories.TaskRepository,
		TimerRepository:           repositories.TimerRepository,
		InformationListRepository: repositories.InformationListRepository,
	}

	services := service.NewService(deps)

	route := gin.Default()
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://gustav.website"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	controller.RouteAPI(route, services)

	srv := server.NewServer(route, server.Port(cfg.Http.Port), server.ReadTimeout(cfg.Http.ReadTimeout),
		server.WriteTimeout(cfg.Http.WriteTimeout), server.ShutdownTimeout(cfg.Http.ShutdownTimeout))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("Run: " + s.String())
	case err := <-srv.Notify():
		log.Println(err, "Run: signal.Notify")
	}

	err = srv.Shutdown()
	if err != nil {
		log.Println(err, "Run: server shutdown")
	}
}
