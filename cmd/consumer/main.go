package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/michelaquino/golang_kafka_example/context"
)

func main() {
	logger := context.GetLogger()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Hack to wait kafka starts. Don't do this in production.
	// See https://docs.docker.com/compose/startup-order/ to better solutions.
	time.Sleep(5000 * time.Millisecond)

	kafkaURL := os.Getenv("KAFKA_URL")
	master, err := sarama.NewConsumer([]string{kafkaURL}, config)
	if err != nil {
		panic(err)
	}

	defer master.Close()

	topic := os.Getenv("KAFKA_TOPIC_TO_CONSUMER")
	partition, err := strconv.Atoi(os.Getenv("KAFKA_PARTITION"))
	if err != nil {
		partition = 0
	}

	consumerPartition, err := master.ConsumePartition(topic, int32(partition), sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	messageCountStart := 0

	go func() {
		for {
			select {
			case err := <-consumerPartition.Errors():
				logger.Error("Consumer", "main", "", "", "Receiving messages", "Error", err.Error())
			case msg := <-consumerPartition.Messages():
				messageCountStart++
				logger.Info("Consumer", "main", "", "", "Receiving messages", "Success", fmt.Sprintf("key: %s - value: %s", msg.Key, msg.Value))
			case <-signals:
				logger.Info("Consumer", "main", "", "", "Receiving messages", "Info", "Interrupt signal is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	logger.Info("Consumer", "main", "", "", "Final", "Info", fmt.Sprintf("Processed %d messages", messageCountStart))
}
