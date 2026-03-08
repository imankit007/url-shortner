package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/imankit007/url-shortner-redirector-service/model"
	"github.com/segmentio/kafka-go"
)

type ClickEventPublisher interface {
	Publish(ctx context.Context, event model.ClickEvent)
}

type kafkaClickEventPublisher struct {
	writer *kafka.Writer
}

func NewClickEventPublisher(writer *kafka.Writer) ClickEventPublisher {
	return &kafkaClickEventPublisher{writer: writer}
}

func (p *kafkaClickEventPublisher) Publish(ctx context.Context, event model.ClickEvent) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal click event: %v", err)
		return
	}

	go func() {
		if err := p.writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(event.ShortCode),
			Value: payload,
		}); err != nil {
			log.Printf("failed to publish click event: %v", err)
		}
	}()
}
