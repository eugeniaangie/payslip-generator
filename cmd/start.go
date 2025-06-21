package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"payslip-generator/controller"
	"payslip-generator/service"
	"payslip-generator/store/postgres_store"
	"payslip-generator/util/config"
	"payslip-generator/util/errs"

	_ "github.com/lib/pq"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func start() {
	const op errs.Op = "main/start"

	// init logger
	var logger = logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Formatter = new(logrus.TextFormatter)                     //default
	logger.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout

	logger.Info("Loading config ...")

	// load environment variables from .env file
	config, err := config.LoadConfig(".")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "LoadConfig",
			"err":   err.Error(),
		}).Error()

		os.Exit(1)
	}

	logger.WithFields(logrus.Fields{
		"op":     op,
		"config": fmt.Sprintf("%+v", config),
	}).Infof("Starting %s service ...", config.App.Name)

	// init store layer
	logger.Info("Creating stores ...")

	postgresStore, err := postgres_store.NewPostgresStore(logger, config.Store.Postgres)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "NewSQLServerStore",
			"err":   err.Error(),
		}).Error()

		os.Exit(1)
	}
	defer postgresStore.Db.Close()

	logger.Info("Postgres store created ...")

	
	// init service layer
	service := service.NewService(
		logger,
		postgresStore,
	)


	// init presentation layer
	controller := controller.NewController(config.App.Timeout, service)

	// init fiber app
	app := fiber.New()

	// CORS middleware configuration
	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}

	app.Use(cors.New(corsConfig))

	// endpoint definitions
	app = defineEndpoints(app, controller)

	// create a channel to signal when the server has stopped
	serverError := make(chan error, 1)

	// start the server
	go func(app *fiber.App, port int, logger *logrus.Logger) {
		if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
			serverError <- err
		}
	}(app, config.App.Port, logger)

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		if err != nil {
			logger.WithError(err).Error("Server stopped unexpectedly")

			os.Exit(1)
		}
	case <-quit:
		logger.Info("Shutdown signal received")
	}

	// Graceful shutdown logic
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
		os.Exit(1)
	}

	logger.Info("Server exiting gracefully")
}
