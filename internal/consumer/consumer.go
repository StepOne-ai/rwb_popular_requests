package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
	"github.com/segmentio/kafka-go"
)

type IUsecase interface {
	ProcessEvent(event *schema.SearchEvent)
}

type Consumer struct {
	reader  *kafka.Reader
	usecase IUsecase
}

func New(brokers []string, topic, group string, uc IUsecase) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  group,
		MinBytes: 1,
		MaxBytes: 1 << 20, // 1MB
	})
	return &Consumer{reader: r, usecase: uc}
}

func (c *Consumer) Run(ctx context.Context) error {
	defer c.reader.Close()
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("consumer: read: %w", err)
		}

		var event schema.SearchEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("consumer: bad message offset=%d: %v", msg.Offset, err)
			continue
		}

		c.usecase.ProcessEvent(&event)
	}
}
