package main

import (
	"context"
	"go.uber.org/zap"
	"ledger-service/routes"
	"ledger-service/services"
	"net/http"
	"os"
	"time"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	services.LoadConfig()
	services.ConnectDB()

	if services.Config.UseRedis {
		services.CheckRedisConnection()
	}

	routes.InitGin()
	router := routes.New()

	server := &http.Server{
		Addr:         services.Config.ServerHost + ":" + services.Config.ServerPort,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("listen", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 15 seconds.
	quit := make(chan os.Signal, 1)
	//signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
