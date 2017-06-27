package main

import (
	"github.com/labstack/echo"
	"github.com/michelaquino/golang_kafka_example/context"
	"github.com/michelaquino/golang_kafka_example/handlers"
)

func main() {
	logger := context.GetLogger()

	echoInstance := echo.New()

	echoInstance.GET("/healthcheck", handlers.Healthcheck)

	logger.Info("Main", "main", "", "", "start app", "success", "Started at port 8888!")
	echoInstance.Logger.Fatal(echoInstance.Start(":8888"))
}
