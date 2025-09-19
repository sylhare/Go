package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/syl/Go/pkg/examples/queue"
)

// ConsumerService represents a service that consumes messages
type ConsumerService struct {
	consumer queue.Consumer
	logger   *log.Logger
}

// NewConsumerService creates a new consumer service
func NewConsumerService(consumer queue.Consumer, logger *log.Logger) *ConsumerService {
	return &ConsumerService{
		consumer: consumer,
		logger:   logger,
	}
}

// Start begins consuming messages from the orders topic
func (cs *ConsumerService) Start(ctx context.Context) error {
	cs.logger.Println("Starting consumer service...")

	if err := cs.consumer.Subscribe(ctx, "orders", cs.handleOrderMessage); err != nil {
		return fmt.Errorf("failed to subscribe to orders topic: %w", err)
	}

	cs.logger.Println("Subscribed to orders topic, waiting for messages...")

	<-ctx.Done()
	cs.logger.Println("Consumer service stopped")
	return ctx.Err()
}

// handleOrderMessage processes an order message
func (cs *ConsumerService) handleOrderMessage(ctx context.Context, message *queue.Message) error {
	cs.logger.Printf("Received message ID: %s from topic: %s", message.ID, message.Topic)

	var order OrderData
	if err := json.Unmarshal(message.Payload, &order); err != nil {
		cs.logger.Printf("Failed to unmarshal order message: %v", err)
		return err
	}

	cs.logger.Printf("Message headers: %v", message.Headers)

	cs.logger.Printf("Processing order: %s", order.OrderID)
	cs.logger.Printf("  Customer ID: %s", order.CustomerID)
	cs.logger.Printf("  Amount: $%.2f", order.Amount)
	cs.logger.Printf("  Created At: %s", order.CreatedAt.Format("2006-01-02 15:04:05"))

	if order.Amount > 100 {
		cs.logger.Printf("High-value order detected: %s (Amount: $%.2f)", order.OrderID, order.Amount)
	}

	cs.logger.Printf("Successfully processed order: %s", order.OrderID)
	return nil
}

// Stop stops the consumer service
func (cs *ConsumerService) Stop() error {
	cs.logger.Println("Stopping consumer service...")
	return cs.consumer.Close()
}
