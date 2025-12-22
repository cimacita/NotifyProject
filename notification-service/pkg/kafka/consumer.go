package kafka

import (
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type IEventHandler interface {
	Handle(msg *kafka.Message) error
}

type Consumer struct {
	consumer *kafka.Consumer
	topic    string
	handler  IEventHandler
}

func NewConsumer(brokers, groupID, topic string, handler IEventHandler) *Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}

	return &Consumer{
		consumer: consumer,
		topic:    topic,
		handler:  handler,
	}
}

func (c *Consumer) Start() {
	err := c.consumer.SubscribeTopics([]string{c.topic}, nil)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	go func() {
		for {
			msg, err := c.consumer.ReadMessage(time.Second)
			if err != nil {
				continue
			}
			log.Printf("Received message: %s\n", string(msg.Value))

			err = c.handler.Handle(msg)
			if err != nil {
				log.Printf("Error handling message: %v", err)
				continue
			}

			_, _ = c.consumer.CommitMessage(msg)
		}
	}()
}

func (c *Consumer) Close() {
	_ = c.consumer.Close()
}
