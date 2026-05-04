package main

import (
	"context"
	"fmt"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/config"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/initializer"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/router"
)

func main() {
	initializer.InitConfig()
	initializer.InitLog()
	logger := initializer.Log
	logger.Info(context.TODO(), "Configuration loaded successfully")
	logger.Info(context.TODO(), "start InitSchema")
	initializer.InitSchema()
	logger.Info(context.TODO(), "start InitPostgresClient")
	initializer.InitPostgresClient()
	logger.Info(context.TODO(), "start InitRedisClient")
	initializer.InitRedisClient()
	logger.Info(context.TODO(), "start router")
	r := router.NewRouter(initializer.GormDB, initializer.Log, initializer.RedisClient)
	logger.Info(context.TODO(), "Core Service starting...")

	serverPort := config.AppConfig.Server.Port
	if err := r.Run(fmt.Sprintf(":%d", serverPort)); err != nil {
		logger.Fatal(context.TODO(), err.Error())
	}
}
