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
	"payslip-generator/util/auth"
	"payslip-generator/util/config"
	"payslip-generator/util/errs"

	"payslip-generator/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @title Payslip Generator API
// @version 1.0
// @description REST API for attendance, overtime, and payroll processing.
// @host localhost:4000
// @BasePath /api
// @schemes http
func start() {
	const op errs.Op = "main/start"
	docs.SwaggerInfo.Title = "Payslip Generator API"
	docs.SwaggerInfo.Description = "API for attendance, overtime, reimbursement, and payroll"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", 4000)
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// init logger
	var logger = logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Formatter = new(logrus.TextFormatter)                     //default
	logger.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout

	logger.Info("Loading config ...")
	config.LoadEnv()

	appConfig, err := config.LoadConfig(".")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "LoadConfig",
			"err":   err.Error(),
		}).Error()
		os.Exit(1)
	}

	auth.SetJWTKey([]byte(config.GetJWTSecret()))

	logger.WithFields(logrus.Fields{
		"op":     op,
		"config": fmt.Sprintf("%+v", appConfig),
	}).Infof("Starting %s service ...", appConfig.App.Name)

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
