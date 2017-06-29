package main

import (
	"github.com/labstack/echo"
	"github.com/michelaquino/golang_kafka_example/context"
	"github.com/michelaquino/golang_kafka_example/handlers"
	apiMiddleware "github.com/michelaquino/golang_kafka_example/middleware"
)

func main() {
	logger := context.GetLogger()

	echoInstance := echo.New()

	// Middlewares
	echoInstance.Use(apiMiddleware.RequestLogDataMiddleware())

	echoInstance.GET("/healthcheck", handlers.Healthcheck)

	// echoInstance.GET("/producer/async", handlers.Healthcheck)
	echoInstance.GET("/producer/sync", handlers.SendSyncMessage)

	logger.Info("Main", "main", "", "", "start app", "success", "Started at port 8888!")
	echoInstance.Logger.Fatal(echoInstance.Start(":8888"))
}
