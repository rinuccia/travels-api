package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rinuccia/travels-api/config"
	"github.com/rinuccia/travels-api/internal/handler"
	"github.com/rinuccia/travels-api/internal/repository/postgres"
	"github.com/rinuccia/travels-api/internal/service"
	"github.com/rinuccia/travels-api/pkg/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// @title Travels API
// @version 1.0
// @description API Server for Travels Application

// @host localhost:8181
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err = godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	//Postgres
	dbPostgres, err := postgres.NewPostgresClient(cfg)
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// Services
	repository := postgres.NewRepository(dbPostgres)
	services := service.NewService(repository)
	handlers := handler.NewHandler(services)

	// Run server
	srv := new(server.Server)
	router := gin.New()
	handlers.InitRoutes(router)

	go func() {
		if err = srv.Run(router, cfg.Port); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	// Waiting signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown
	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err = dbPostgres.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
