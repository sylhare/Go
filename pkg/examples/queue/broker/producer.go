// Package broker provides Producer and Consumer implementations
package broker

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/syl/Go/pkg/examples/queue"
)

// QueueProducer implements the Producer interface using a Queue
type QueueProducer struct {
	queue queue.Queue
}

// NewQueueProducer creates a new producer that uses the provided queue
func NewQueueProducer(q queue.Queue) *QueueProducer {
	return &QueueProducer{
		queue: q,
	}
}

// Publish sends a message to the specified topic
func (p *QueueProducer) Publish(ctx context.Context, topic string, payload []byte, headers map[string]string) error {
	message := &queue.Message{
		ID:        uuid.New().String(),
		Topic:     topic,
		Payload:   payload,
		Headers:   headers,
		Timestamp: time.Now(),
	}

	return p.queue.Enqueue(ctx, topic, message)
}

// Close closes the producer
func (p *QueueProducer) Close() error {
	return nil
}
