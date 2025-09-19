package broker

import (
	"context"
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
// It accepts any queue.Queue implementation, making broker tests implementation-agnostic
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
