package infrastructure

import "github.com/segmentio/kafka-go"

const ClickEventsTopic = "click-events"

func NewKafkaWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    ClickEventsTopic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
		AllowAutoTopicCreation: true,
	}
}
