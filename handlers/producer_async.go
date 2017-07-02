package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
	"github.com/michelaquino/golang_kafka_example/context"
	apiMiddleware "github.com/michelaquino/golang_kafka_example/middleware"
	"github.com/michelaquino/golang_kafka_example/models"
)

// SendAsyncMessage is a method that send a message to kafka
func SendAsyncMessage(echoContext echo.Context) error {
	logger := context.GetLogger()
	requestLogData := echoContext.Get(apiMiddleware.RequestIDKey).(models.RequestLogData)

	type messageToEnqueue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var message messageToEnqueue
	if err := echoContext.Bind(&message); err != nil {
		logger.Error("Handlers", "SendAsyncMessage", requestLogData.ID, requestLogData.OriginIP, "Bind payload to object", "Error", err.Error())
		return echoContext.NoContent(http.StatusBadGateway)
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	kafkaURL := os.Getenv("KAFKA_URL")

	producer, err := sarama.NewAsyncProducer([]string{kafkaURL}, config)
	if err != nil {
		logger.Error("Handlers", "SendAsyncMessage", requestLogData.ID, requestLogData.OriginIP, "Try to connect to kafka host", "Error", err.Error())
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	defer producer.Close()

	asyncTopic := os.Getenv("KAFKA_ASYNC_TOPIC")
	partitionToEnquene, err := strconv.Atoi(os.Getenv("KAFKA_PARTITION"))
	if err != nil {
		partitionToEnquene = 0
	}

	kafkaMessage := sarama.ProducerMessage{
		Topic:     asyncTopic,
		Key:       sarama.StringEncoder(message.Key),
		Value:     sarama.StringEncoder(message.Value),
		Partition: int32(partitionToEnquene),
	}

	for {
		select {
		case producer.Input() <- &kafkaMessage:
			logger.Info("Handlers", "SendAsyncMessage", requestLogData.ID, requestLogData.OriginIP, "Send async message", "Success", "")
			return echoContext.String(http.StatusOK, fmt.Sprintf("Message enqueued"))

		case err := <-producer.Errors():
			logger.Error("Handlers", "SendAsyncMessage", requestLogData.ID, requestLogData.OriginIP, "Send async message", "Error", err.Error())
			return echoContext.NoContent(http.StatusInternalServerError)
		}
	}
}
