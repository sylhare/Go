// Package example provides working examples of the queue system
package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/syl/Go/pkg/examples/queue"
)

// OrderData represents an example order message
type OrderData struct {
	OrderID    string    `json:"order_id"`
	CustomerID string    `json:"customer_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

// ProducerService represents a service that produces messages
type ProducerService struct {
	producer queue.Producer
	logger   *log.Logger
}

// NewProducerService creates a new producer service
func NewProducerService(producer queue.Producer, logger *log.Logger) *ProducerService {
	return &ProducerService{
		producer: producer,
		logger:   logger,
	}
}

// Start begins producing messages periodically
func (ps *ProducerService) Start(ctx context.Context) error {
	ps.logger.Println("Starting producer service...")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	orderID := 1

	for {
		select {
		case <-ctx.Done():
			ps.logger.Println("Producer service stopped")
			return ctx.Err()
		case <-ticker.C:
			order := OrderData{
				OrderID:    fmt.Sprintf("order-%d", orderID),
				CustomerID: fmt.Sprintf("customer-%d", (orderID%5)+1),
				Amount:     float64(orderID * 10),
				CreatedAt:  time.Now(),
			}

			payload, err := json.Marshal(order)
			if err != nil {
				ps.logger.Printf("Failed to marshal order: %v", err)
				continue
			}

			headers := map[string]string{
				"source":      "producer-service",
				"message_type": "order",
				"version":     "1.0",
			}

			if err := ps.producer.Publish(ctx, "orders", payload, headers); err != nil {
				ps.logger.Printf("Failed to publish order %s: %v", order.OrderID, err)
			} else {
				ps.logger.Printf("Published order: %s (Customer: %s, Amount: %.2f)",
					order.OrderID, order.CustomerID, order.Amount)
			}

			orderID++
		}
	}
}

// Stop stops the producer service
func (ps *ProducerService) Stop() error {
	ps.logger.Println("Stopping producer service...")
	return ps.producer.Close()
}
