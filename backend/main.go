package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sfu-teamproject/smartbuy/backend/app"
	"github.com/sfu-teamproject/smartbuy/backend/logger"
	"github.com/sfu-teamproject/smartbuy/backend/storage/postgres"
)

// @title           Smartbuy API
// @version         1.0
// @description     API Server for Smartbuy application
// @termsOfService  http://swagger.io/terms/
// @openapi 3.0.0
// @host            localhost:8081
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	godotenv.Load()
	logger, err := logger.NewConsoleLogger()
	if err != nil {
		log.Fatalf("Error getting logger: %v", err)
	}
	server := &http.Server{
		Addr:     ":8081",
		ErrorLog: logger.Error,
	}
	postgres, err := postgres.NewPostgresDB(false)
	if err != nil {
		logger.Errorf("Error creating database: %v", err)
		os.Exit(1)
	}
	a := app.NewApp(logger, server, postgres)
	a.Server.Handler = a.NewRouter()
	a.Log.Infof("Starting server on %s", a.Server.Addr)
	err = a.Server.ListenAndServe()
	if err != nil {
		a.Log.Errorf("Error starting server: %v", err)
		os.Exit(1)
	}
}
