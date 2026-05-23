package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func New(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, event *schema.SearchEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("producer: marshal: %w", err)
	}
	return p.writer.WriteMessages(ctx, kafka.Message{Value: body})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
