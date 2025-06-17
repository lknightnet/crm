package internal

import (
	"github.com/gin-gonic/gin"
	"jwt-service/config"
	"jwt-service/internal/controller"
	"jwt-service/internal/model"
	"jwt-service/internal/repository"
	"jwt-service/internal/service"
	"jwt-service/pkg/database"
	"jwt-service/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	connection, err := database.NewConnection(cfg.Database.URI, 0, 5, &model.RefreshToken{}, &model.AccessToken{})
	if err != nil {
		panic(err)
	}

	repositories := repository.NewRepositories(connection)

	deps := &service.DependenciesService{
		SignKey:       []byte(cfg.JWT.SignKey),
		AccessExpiry:  cfg.JWT.AccessExpiry,
		RefreshExpiry: cfg.JWT.RefreshExpiry,
		JWTRepository: repositories.JWTRepository,
	}

	services := service.NewService(deps)

	route := gin.New()
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
