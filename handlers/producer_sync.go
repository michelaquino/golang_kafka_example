package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
	"github.com/michelaquino/golang_kafka_example/context"
	apiMiddleware "github.com/michelaquino/golang_kafka_example/middleware"
	"github.com/michelaquino/golang_kafka_example/models"
)

// SendSyncMessage is a method that send a message to kafka
func SendSyncMessage(echoContext echo.Context) error {
	logger := context.GetLogger()
	requestLogData := echoContext.Get(apiMiddleware.RequestIDKey).(models.RequestLogData)

	type messageToEnqueue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var message messageToEnqueue
	if err := echoContext.Bind(&message); err != nil {
		logger.Error("Handlers", "SendSyncMessage", requestLogData.ID, requestLogData.OriginIP, "Bind payload to object", "Error", err.Error())
		return echoContext.NoContent(http.StatusBadGateway)
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	kafkaURL := os.Getenv("KAFKA_URL")
	producer, err := sarama.NewSyncProducer([]string{kafkaURL}, config)
	if err != nil {
		logger.Error("Handlers", "SendSyncMessage", requestLogData.ID, requestLogData.OriginIP, "Try to connect to kafka host", "Error", err.Error())
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	defer producer.Close()

	asyncTopic := os.Getenv("KAFKA_SYNC_TOPIC")
	kafkaMessage := sarama.ProducerMessage{
		Topic: asyncTopic,
		Key:   sarama.StringEncoder(message.Key),
		Value: sarama.StringEncoder(message.Value),
	}

	partition, offset, err := producer.SendMessage(&kafkaMessage)
	if err != nil {
		logger.Error("Handlers", "SendSyncMessage", requestLogData.ID, requestLogData.OriginIP, "Send sync message", "Error", err.Error())
		return echoContext.String(http.StatusInternalServerError, "Error to send message")
	}

	logger.Info("Handlers", "SendSyncMessage", requestLogData.ID, requestLogData.OriginIP, "Send sync message", "Success", "")
	return echoContext.String(http.StatusOK, fmt.Sprintf("Message stored in topic %s, partition %d, offset %d", asyncTopic, partition, offset))
}
