package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(brokers string) *Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"acks":              "all",
	})
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			case kafka.Error:
				log.Printf("Producer error: %v\n", ev.Error())
			}
		}
	}()

	return &Producer{producer: p}
}

func (p *Producer) Produce(topic, key string, message []byte) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: message,
	}

	return p.producer.Produce(msg, nil)
}

func (p *Producer) Close() {
	p.producer.Flush(10000)
	p.producer.Close()
}
