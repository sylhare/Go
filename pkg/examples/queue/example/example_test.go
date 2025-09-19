package example

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syl/Go/pkg/examples/queue"
	"github.com/syl/Go/pkg/examples/queue/broker"
	"github.com/syl/Go/pkg/examples/queue/inmemory"
	"github.com/syl/Go/pkg/examples/queue/inmemory/testutils"
)

// Test helpers for order-specific functionality

func createTestOrder(id string) OrderData {
	return OrderData{
		OrderID:    "test-order-" + id,
		CustomerID: "test-customer-" + id,
		Amount:     99.99,
		CreatedAt:  time.Now(),
	}
}

func createTestMessageFromOrder(order OrderData) (*queue.Message, error) {
	payload, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}
	
	return &queue.Message{
		ID:      "test-msg-" + order.OrderID,
		Topic:   "orders",
		Payload: payload,
		Headers: map[string]string{
			"source":       "test",
			"message_type": "order",
		},
		Timestamp: time.Now(),
	}, nil
}

func assertOrderMessage(t *testing.T, msg *queue.Message) {
	t.Helper()
	
	var order OrderData
	err := json.Unmarshal(msg.Payload, &order)
	require.NoError(t, err, "Should unmarshal order data")
	
	assert.NotEmpty(t, order.OrderID, "OrderID should be set")
	assert.NotEmpty(t, order.CustomerID, "CustomerID should be set")
	assert.Greater(t, order.Amount, float64(0), "Amount should be positive")
	
	expectedHeaders := map[string]string{
		"source":       "producer-service",
		"message_type": "order",
		"version":      "1.0",
	}
	
	for key, expectedValue := range expectedHeaders {
		assert.Equal(t, expectedValue, msg.Headers[key], "Header %s should match", key)
	}
}

func TestExampleServices(t *testing.T) {
	t.Run("ProducerService", func(t *testing.T) {
		q := inmemory.NewInMemoryQueue()
		producer := broker.NewQueueProducer(q)
		consumer := broker.NewQueueConsumer(q)
		fixture := testutils.NewServiceFixture(t, q, producer, consumer)
		
		producerLogger := testutils.CreateLogger("TEST-PRODUCER-PRODUCER")
		producerService := NewProducerService(fixture.Producer, producerLogger)
		t.Cleanup(func() { producerService.Stop() })
		
		t.Run("GeneratesMessages", func(t *testing.T) {
			fixture.RunServiceWithTimeout(producerService.Start, 3*time.Second)

			size, err := fixture.Queue.Size(context.Background(), "orders")
			require.NoError(t, err)
			assert.Greater(t, size, 0, "Expected producer to generate messages")
			t.Logf("Producer generated %d messages", size)
		})
		
		t.Run("ValidatesMessageFormatAndHeaders", func(t *testing.T) {
			msg, err := fixture.Queue.Dequeue(context.Background(), "orders")
			require.NoError(t, err)
			require.NotNil(t, msg, "Expected message to exist")

			assertOrderMessage(t, msg)
		})
	})

	t.Run("ConsumerService", func(t *testing.T) {
		q := inmemory.NewInMemoryQueue()
		producer := broker.NewQueueProducer(q)
		consumer := broker.NewQueueConsumer(q)
		fixture := testutils.NewServiceFixture(t, q, producer, consumer)
		
		consumerLogger := testutils.CreateLogger("TEST-CONSUMER-CONSUMER")
		consumerService := NewConsumerService(fixture.Consumer, consumerLogger)
		t.Cleanup(func() { consumerService.Stop() })
		
		testOrder := createTestOrder("123")
		
		t.Run("EnqueueTestMessage", func(t *testing.T) {
			testMsg, err := createTestMessageFromOrder(testOrder)
			require.NoError(t, err)

			err = fixture.Queue.Enqueue(context.Background(), "orders", testMsg)
			require.NoError(t, err)
			
			fixture.AssertQueueSize("orders", 1, "Should have one message in queue")
		})
		
		t.Run("ConsumesMessage", func(t *testing.T) {
			fixture.RunServiceWithTimeout(consumerService.Start, 2*time.Second)

			size, err := fixture.Queue.Size(context.Background(), "orders")
			require.NoError(t, err)
			t.Logf("Queue size after consumer ran: %d", size)
		})
	})

	t.Run("Integration", func(t *testing.T) {
		t.Run("EndToEndFlow", func(t *testing.T) {
			q := inmemory.NewInMemoryQueue()
			producer := broker.NewQueueProducer(q)
			consumer := broker.NewQueueConsumer(q)
			fixture := testutils.NewServiceFixture(t, q, producer, consumer)
			
			processedMessages := make(chan string, 10)
			ctx := context.Background()
			
			handler := func(ctx context.Context, message *queue.Message) error {
				var order OrderData
				err := json.Unmarshal(message.Payload, &order)
				require.NoError(t, err)

				processedMessages <- order.OrderID
				return nil
			}

			err := fixture.Consumer.Subscribe(ctx, "orders", handler)
			require.NoError(t, err, "Should subscribe successfully")

			testOrder := createTestOrder("integration")
			payload, err := json.Marshal(testOrder)
			require.NoError(t, err)

			err = fixture.Producer.Publish(ctx, "orders", payload, map[string]string{
				"source":       "integration-test",
				"message_type": "order",
			})
			require.NoError(t, err)

			select {
			case orderID := <-processedMessages:
				assert.Equal(t, testOrder.OrderID, orderID, "Should process the correct order")
				t.Logf("Successfully processed order: %s", orderID)
			case <-time.After(3 * time.Second):
				t.Fatal("Timeout waiting for message processing")
			}
		})
	})
}
