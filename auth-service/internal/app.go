package internal

import (
	"auth-service/config"
	"auth-service/internal/controller"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/database"
	"auth-service/pkg/server"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	connection, err := database.NewConnection(cfg.Database.URI, 0, 5, &model.RouteLink{})
	if err != nil {
		panic(err)
	}

	repositories := repository.NewRepositories(connection)

	deps := &service.DependenciesService{
		AuthSignature:  cfg.Auth.Signature,
		JWTServiceApi:  cfg.JWTService.Api,
		JWTServicePort: cfg.JWTService.Port,
		UserRepository: repositories.UserRepository,
		AuthRepository: repositories.AuthRepository,
	}

	services := service.NewService(deps)

	route := gin.Default()
	controller.RouteAPI(route, services)
	//route.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"http://localhost:3000"},
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	//	AllowHeaders:     []string{"Authorization", "Content-Type"},
	//	AllowCredentials: true,
	//}))
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
