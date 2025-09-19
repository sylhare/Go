package broker

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syl/Go/pkg/examples/queue"
)

const (
	// ConsumerStartupDelay allows consumers to start processing
	ConsumerStartupDelay = 200 * time.Millisecond
	// DefaultTestTimeout for broker test assertions
	DefaultTestTimeout = 5 * time.Second
)

// BrokerTestFixture provides test infrastructure for broker testing
type BrokerTestFixture struct {
	Queue    queue.Queue
	Producer queue.Producer
	Consumer queue.Consumer
	Ctx      context.Context
	T        *testing.T
}

// NewBrokerTestFixture creates a fixture for testing broker components
func NewBrokerTestFixture(t *testing.T, q queue.Queue) *BrokerTestFixture {
	t.Helper()

	producer := NewQueueProducer(q)
	consumer := NewQueueConsumer(q)
	ctx := context.Background()

	t.Cleanup(func() {
		consumer.Close()
		producer.Close()
		q.Close()
	})

	return &BrokerTestFixture{
		Queue:    q,
		Producer: producer,
		Consumer: consumer,
		Ctx:      ctx,
		T:        t,
	}
}

// AssertQueueSize verifies the queue size for a topic
func (f *BrokerTestFixture) AssertQueueSize(topic string, expectedSize int, msgAndArgs ...interface{}) {
	f.T.Helper()

	size, err := f.Queue.Size(f.Ctx, topic)
	require.NoError(f.T, err, "Should get queue size")
	assert.Equal(f.T, expectedSize, size, msgAndArgs...)
}

// PublishMessages publishes multiple test messages to a topic
func (f *BrokerTestFixture) PublishMessages(topic string, messages []string) {
	f.T.Helper()

	for _, msg := range messages {
		err := f.Producer.Publish(f.Ctx, topic, []byte(msg), nil)
		require.NoError(f.T, err, "Should publish message %s", msg)
	}
}

// AssertMessagesReceived waits for and validates received messages
func (f *BrokerTestFixture) AssertMessagesReceived(msgChan <-chan *queue.Message, expectedCount int, timeout time.Duration) {
	f.T.Helper()

	receivedCount := 0
	timeoutChan := time.After(timeout)

	for receivedCount < expectedCount {
		select {
		case msg := <-msgChan:
			f.T.Logf("Received message: %s", string(msg.Payload))
			receivedCount++
		case <-timeoutChan:
			f.T.Fatalf("Timeout waiting for messages, received %d of %d", receivedCount, expectedCount)
		}
	}

	assert.Equal(f.T, expectedCount, receivedCount, "Should receive expected number of messages")
}

func TestBroker(t *testing.T) {
	t.Run("Producer", func(t *testing.T) {
		t.Run("PublishMessage", func(t *testing.T) {
			q := queue.NewMock()
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
			q := queue.NewMock()
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
			q := queue.NewMock()
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
			q := queue.NewMock()
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
