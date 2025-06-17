package internal

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	"user-service/config"
	"user-service/internal/controller"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg/database"
	"user-service/pkg/server"
)

func Run(cfg *config.Config) {
	connection, err := database.NewConnection(cfg.Database.URI, 0, 5, &model.User{})
	if err != nil {
		panic(err)
	}

	repositories := repository.NewRepositories(connection)

	deps := &service.DependenciesService{
		JWTServiceApi:  cfg.JWTService.Api,
		JWTServicePort: cfg.JWTService.Port,
		UserRepository: repositories.UserRepository,
	}

	services := service.NewService(deps)

	route := gin.Default()
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://gustav.website", "http://crm-project-service-1:8012", "http://localhost:3000"},
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
