package broker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syl/Go/pkg/examples/queue"
)

// fakeQueue is a simple test-only queue implementation for broker testing
// This allows broker tests to be completely independent of any specific queue implementation
type fakeQueue struct {
	topics map[string]chan *queue.Message
	closed bool
	mutex  sync.RWMutex
}

// newFakeQueue creates a new fake queue for testing
func newFakeQueue() *fakeQueue {
	return &fakeQueue{
		topics: make(map[string]chan *queue.Message),
		closed: false,
	}
}

func (q *fakeQueue) Enqueue(ctx context.Context, topic string, message *queue.Message) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.closed {
		return errors.New("queue is closed")
	}

	if _, exists := q.topics[topic]; !exists {
		q.topics[topic] = make(chan *queue.Message, 100) // Buffered for testing
	}

	select {
	case q.topics[topic] <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *fakeQueue) Dequeue(ctx context.Context, topic string) (*queue.Message, error) {
	q.mutex.RLock()
	ch, exists := q.topics[topic]
	closed := q.closed
	q.mutex.RUnlock()

	if closed {
		return nil, errors.New("queue is closed")
	}

	if !exists {
		return nil, nil // No messages in topic
	}

	select {
	case msg := <-ch:
		return msg, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, nil // No messages available
	}
}

func (q *fakeQueue) Size(ctx context.Context, topic string) (int, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if q.closed {
		return 0, errors.New("queue is closed")
	}

	if ch, exists := q.topics[topic]; exists {
		return len(ch), nil
	}
	return 0, nil
}

func (q *fakeQueue) Topics(ctx context.Context) ([]string, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if q.closed {
		return nil, errors.New("queue is closed")
	}

	topics := make([]string, 0, len(q.topics))
	for topic := range q.topics {
		topics = append(topics, topic)
	}
	return topics, nil
}

func (q *fakeQueue) Close() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.closed {
		return nil
	}

	q.closed = true
	for _, ch := range q.topics {
		close(ch)
	}
	return nil
}

func TestBroker(t *testing.T) {
	t.Run("Producer", func(t *testing.T) {
		t.Run("PublishMessage", func(t *testing.T) {
			q := newFakeQueue()
			fixture := NewBrokerTestFixture(t, q)
			topic := "test-topic"
			payload := []byte("test message")
			headers := map[string]string{"key": "value"}

			err := fixture.Producer.Publish(fixture.Ctx, topic, payload, headers)
			require.NoError(t, err, "Should publish message successfully")

			fixture.AssertQueueSize(topic, 1, "Queue should contain exactly one message")

			msg, err := fixture.Queue.Dequeue(fixture.Ctx, topic)
			require.NoError(t, err, "Should dequeue message")
			require.NotNil(t, msg, "Message should not be nil")

			assert.Equal(t, payload, msg.Payload, "Payload should match")
			assert.Equal(t, "value", msg.Headers["key"], "Header should match")
			assert.Equal(t, topic, msg.Topic, "Topic should match")
		})
	})

	t.Run("Consumer", func(t *testing.T) {
		t.Run("Subscribe", func(t *testing.T) {
			q := newFakeQueue()
			fixture := NewBrokerTestFixture(t, q)
			topic := "test-topic"
			receivedMessages := make(chan *queue.Message, 10)

			handler := func(ctx context.Context, message *queue.Message) error {
				receivedMessages <- message
				return nil
			}

			err := fixture.Consumer.Subscribe(fixture.Ctx, topic, handler)
			require.NoError(t, err, "Should subscribe successfully")

			time.Sleep(ConsumerStartupDelay)

			testMessages := []string{"message1", "message2", "message3"}
			fixture.PublishMessages(topic, testMessages)

			fixture.AssertMessagesReceived(receivedMessages, len(testMessages), DefaultTestTimeout)

			err = fixture.Consumer.Unsubscribe(fixture.Ctx, topic)
			require.NoError(t, err, "Should unsubscribe successfully")
		})

		t.Run("MultipleSubscriptions", func(t *testing.T) {
			q := newFakeQueue()
			fixture := NewBrokerTestFixture(t, q)
			topic1, topic2 := "topic1", "topic2"

			topic1Messages := make(chan *queue.Message, 10)
			topic2Messages := make(chan *queue.Message, 10)

			handler1 := func(ctx context.Context, message *queue.Message) error {
				topic1Messages <- message
				return nil
			}

			handler2 := func(ctx context.Context, message *queue.Message) error {
				topic2Messages <- message
				return nil
			}

			err := fixture.Consumer.Subscribe(fixture.Ctx, topic1, handler1)
			require.NoError(t, err, "Should subscribe to topic1")

			err = fixture.Consumer.Subscribe(fixture.Ctx, topic2, handler2)
			require.NoError(t, err, "Should subscribe to topic2")

			time.Sleep(ConsumerStartupDelay)

			err = fixture.Producer.Publish(fixture.Ctx, topic1, []byte("message for topic1"), nil)
			require.NoError(t, err, "Should publish to topic1")

			err = fixture.Producer.Publish(fixture.Ctx, topic2, []byte("message for topic2"), nil)
			require.NoError(t, err, "Should publish to topic2")

			timeout := time.After(DefaultTestTimeout)
			receivedTopic1, receivedTopic2 := false, false

			for !receivedTopic1 || !receivedTopic2 {
				select {
				case msg := <-topic1Messages:
					t.Logf("Received from topic1: %s", string(msg.Payload))
					receivedTopic1 = true
				case msg := <-topic2Messages:
					t.Logf("Received from topic2: %s", string(msg.Payload))
					receivedTopic2 = true
				case <-timeout:
					t.Fatalf("Timeout waiting for messages. Topic1: %v, Topic2: %v", receivedTopic1, receivedTopic2)
				}
			}

			assert.True(t, receivedTopic1, "Should receive message from topic1")
			assert.True(t, receivedTopic2, "Should receive message from topic2")
		})

		t.Run("ConcurrentConsumers", func(t *testing.T) {
			q := newFakeQueue()
			fixture := NewBrokerTestFixture(t, q)
			topic := "concurrent-topic"
			numConsumers := 3
			numMessages := 10

			allMessages := make(chan *queue.Message, numMessages)
			var wg sync.WaitGroup

			consumers := make([]*QueueConsumer, numConsumers)
			for i := 0; i < numConsumers; i++ {
				consumers[i] = NewQueueConsumer(fixture.Queue)
				consumerRef := consumers[i]
				t.Cleanup(func() { consumerRef.Close() })
				
				handler := func(ctx context.Context, message *queue.Message) error {
					allMessages <- message
					return nil
				}

				err := consumers[i].Subscribe(fixture.Ctx, topic, handler)
				require.NoError(t, err, "Should subscribe consumer %d", i)
			}

			time.Sleep(ConsumerStartupDelay)

			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < numMessages; i++ {
					err := fixture.Producer.Publish(fixture.Ctx, topic, []byte(fmt.Sprintf("message-%d", i)), nil)
					assert.NoError(t, err, "Should publish message %d", i)
				}
			}()

			wg.Wait()

			fixture.AssertMessagesReceived(allMessages, numMessages, 10*time.Second)
		})
	})
}
