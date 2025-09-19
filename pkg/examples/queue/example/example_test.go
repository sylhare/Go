package example

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syl/Go/pkg/examples/queue"
	"github.com/syl/Go/pkg/examples/queue/broker"
	"github.com/syl/Go/pkg/examples/queue/inmemory"
)

func TestExampleServices(t *testing.T) {
	t.Run("ProducerService", func(t *testing.T) {
		var q *inmemory.InMemoryQueue
		var producer *broker.QueueProducer
		var service *ProducerService
		
		t.Run("Setup", func(t *testing.T) {
			q = inmemory.NewInMemoryQueue()
			require.NotNil(t, q)
			
			producer = broker.NewQueueProducer(q)
			require.NotNil(t, producer)
			
			logger := log.New(os.Stdout, "[TEST-PRODUCER] ", log.LstdFlags)
			service = NewProducerService(producer, logger)
			require.NotNil(t, service)
		})
		
		t.Run("GeneratesMessages", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				err := service.Start(ctx)
				assert.True(t, err == nil || err == context.DeadlineExceeded)
			}()

			wg.Wait()

			size, err := q.Size(context.Background(), "orders")
			require.NoError(t, err)
			assert.Greater(t, size, 0, "Expected producer to generate messages")
			t.Logf("Producer generated %d messages", size)
		})
		
		t.Run("ValidatesMessageFormatAndHeaders", func(t *testing.T) {
			msg, err := q.Dequeue(context.Background(), "orders")
			require.NoError(t, err)
			require.NotNil(t, msg, "Expected message to exist")

			var order OrderData
			err = json.Unmarshal(msg.Payload, &order)
			require.NoError(t, err, "Should unmarshal order data")
			
			assert.NotEmpty(t, order.OrderID, "OrderID should be set")
			assert.NotEmpty(t, order.CustomerID, "CustomerID should be set")
			assert.Greater(t, order.Amount, float64(0), "Amount should be positive")
			
			assert.Equal(t, "producer-service", msg.Headers["source"])
			assert.Equal(t, "order", msg.Headers["message_type"])
			assert.Equal(t, "1.0", msg.Headers["version"])
		})
		
		t.Cleanup(func() {
			if producer != nil {
				producer.Close()
			}
			if q != nil {
				q.Close()
			}
		})
	})

	t.Run("ConsumerService", func(t *testing.T) {
		var q *inmemory.InMemoryQueue
		var consumer *broker.QueueConsumer
		var service *ConsumerService
		var testOrder OrderData
		
		t.Run("Setup", func(t *testing.T) {
			q = inmemory.NewInMemoryQueue()
			require.NotNil(t, q)

			consumer = broker.NewQueueConsumer(q)
			require.NotNil(t, consumer)

			logger := log.New(os.Stdout, "[TEST-CONSUMER] ", log.LstdFlags)
			service = NewConsumerService(consumer, logger)
			require.NotNil(t, service)

			testOrder = OrderData{
				OrderID:    "test-order-123",
				CustomerID: "test-customer-456",
				Amount:     99.99,
				CreatedAt:  time.Now(),
			}
		})
		
		t.Run("EnqueueTestMessage", func(t *testing.T) {
			payload, err := json.Marshal(testOrder)
			require.NoError(t, err)

			testMsg := &queue.Message{
				ID:      "test-msg",
				Topic:   "orders",
				Payload: payload,
				Headers: map[string]string{
					"source":       "test",
					"message_type": "order",
				},
				Timestamp: time.Now(),
			}

			err = q.Enqueue(context.Background(), "orders", testMsg)
			require.NoError(t, err)
			
			size, err := q.Size(context.Background(), "orders")
			require.NoError(t, err)
			assert.Equal(t, 1, size, "Should have one message in queue")
		})
		
		t.Run("ConsumesMessage", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				err := service.Start(ctx)
				assert.True(t, err == nil || err == context.DeadlineExceeded)
			}()

			wg.Wait()

			size, err := q.Size(context.Background(), "orders")
			require.NoError(t, err)
			t.Logf("Queue size after consumer ran: %d", size)
		})
		
		t.Cleanup(func() {
			if consumer != nil {
				consumer.Close()
			}
			if q != nil {
				q.Close()
			}
		})
	})

	t.Run("Integration", func(t *testing.T) {
		t.Run("EndToEndFlow", func(t *testing.T) {
			q := inmemory.NewInMemoryQueue()
			defer q.Close()

			producer := broker.NewQueueProducer(q)
			defer producer.Close()

			consumer := broker.NewQueueConsumer(q)
			defer consumer.Close()

			processedMessages := make(chan string, 10)
			ctx := context.Background()
			
			handler := func(ctx context.Context, message *queue.Message) error {
				var order OrderData
				err := json.Unmarshal(message.Payload, &order)
				require.NoError(t, err)

				processedMessages <- order.OrderID
				return nil
			}

			err := consumer.Subscribe(ctx, "orders", handler)
			require.NoError(t, err, "Should subscribe successfully")

			testOrder := OrderData{
				OrderID:    "integration-test-order",
				CustomerID: "integration-test-customer",
				Amount:     123.45,
				CreatedAt:  time.Now(),
			}

			payload, err := json.Marshal(testOrder)
			require.NoError(t, err)

			err = producer.Publish(ctx, "orders", payload, map[string]string{
				"source":       "integration-test",
				"message_type": "order",
			})
			require.NoError(t, err)

			select {
			case orderID := <-processedMessages:
				assert.Equal(t, "integration-test-order", orderID, "Should process the correct order")
				t.Logf("Successfully processed order: %s", orderID)
			case <-time.After(3 * time.Second):
				t.Fatal("Timeout waiting for message processing")
			}
		})
	})
}
