// Package kafka hiện thực app.Publisher, gửi domain event dạng JSON lên Kafka.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/my-org/base-go/internal/order/domain"
)

type Publisher struct {
	w *kafka.Writer
}

func NewPublisher(brokers []string, topic string) *Publisher {
	return &Publisher{w: &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 10 * time.Millisecond,
	}}
}

type message struct {
	Type    string `json:"type"`
	OrderID string `json:"order_id"`
}

func (p *Publisher) Publish(ctx context.Context, events []any) error {
	msgs := make([]kafka.Message, 0, len(events))
	for _, e := range events {
		switch ev := e.(type) {
		case domain.OrderPlaced:
			b, err := json.Marshal(message{Type: "OrderPlaced", OrderID: ev.OrderID.String()})
			if err != nil {
				return err
			}
			msgs = append(msgs, kafka.Message{Key: []byte(ev.OrderID.String()), Value: b})
		default:
			return fmt.Errorf("kafka: unknown event %T", e)
		}
	}
	if len(msgs) == 0 {
		return nil
	}
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *Publisher) Close() error { return p.w.Close() }
