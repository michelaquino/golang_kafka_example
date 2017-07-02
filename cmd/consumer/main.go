package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	"github.com/michelaquino/golang_kafka_example/context"
)

func main() {
	logger := context.GetLogger()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	time.Sleep(5000 * time.Millisecond) // Hack to wait kafka starts. Don't do this in production

	kafkaURL := os.Getenv("KAFKA_URL")
	master, err := sarama.NewConsumer([]string{kafkaURL}, config)
	if err != nil {
		panic(err)
	}
	defer master.Close()

	topic := os.Getenv("KAFKA_TOPIC_TO_CONSUMER")
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
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
			case err := <-consumer.Errors():
				logger.Error("Consumer", "main", "", "", "Receiving messages", "Error", err.Error())
			case msg := <-consumer.Messages():
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
