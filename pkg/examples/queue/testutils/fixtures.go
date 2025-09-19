package testutils

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syl/Go/pkg/examples/queue"
)

// BaseFixture provides common test infrastructure
type BaseFixture struct {
	Queue queue.Queue
	Ctx   context.Context
	T     *testing.T
}

// NewBaseFixture creates a base fixture with queue and context
func NewBaseFixture(t *testing.T, queue queue.Queue) *BaseFixture {
	t.Helper()
	
	ctx := context.Background()
	
	t.Cleanup(func() {
		queue.Close()
	})
	
	return &BaseFixture{
		Queue: queue,
		Ctx:   ctx,
		T:     t,
	}
}

// BrokerFixture extends BaseFixture with producer and consumer
type BrokerFixture struct {
	*BaseFixture
	Producer queue.Producer
	Consumer queue.Consumer
}

// NewBrokerFixture creates a fixture with producer and consumer
func NewBrokerFixture(t *testing.T, queue queue.Queue, producer queue.Producer, consumer queue.Consumer) *BrokerFixture {
	t.Helper()
	
	base := NewBaseFixture(t, queue)
	
	t.Cleanup(func() {
		consumer.Close()
		producer.Close()
	})
	
	return &BrokerFixture{
		BaseFixture: base,
		Producer:    producer,
		Consumer:    consumer,
	}
}

// CreateMessage creates a test message with default values
func (f *BaseFixture) CreateMessage(id, topic string, payload []byte) *queue.Message {
	f.T.Helper()
	
	return &queue.Message{
		ID:        id,
		Topic:     topic,
		Payload:   payload,
		Headers:   map[string]string{"key": "value"},
		Timestamp: time.Now(),
	}
}

// CreateMessageWithHeaders creates a test message with custom headers
func (f *BaseFixture) CreateMessageWithHeaders(id, topic string, payload []byte, headers map[string]string) *queue.Message {
	f.T.Helper()
	
	return &queue.Message{
		ID:        id,
		Topic:     topic,
		Payload:   payload,
		Headers:   headers,
		Timestamp: time.Now(),
	}
}

// AssertQueueSize verifies the queue size for a topic
func (f *BaseFixture) AssertQueueSize(topic string, expectedSize int, message string) {
	f.T.Helper()
	
	size, err := f.Queue.Size(f.Ctx, topic)
	require.NoError(f.T, err, "Should get queue size")
	assert.Equal(f.T, expectedSize, size, message)
}

// AssertTopicsContain verifies that topics contain expected values
func (f *BaseFixture) AssertTopicsContain(expectedTopics ...string) {
	f.T.Helper()
	
	topics, err := f.Queue.Topics(f.Ctx)
	require.NoError(f.T, err, "Should get topics successfully")
	
	topicMap := make(map[string]bool)
	for _, topic := range topics {
		topicMap[topic] = true
	}
	
	for _, expected := range expectedTopics {
		assert.True(f.T, topicMap[expected], "Should contain topic %s", expected)
	}
}

// PublishMessages publishes multiple test messages to a topic
func (f *BrokerFixture) PublishMessages(topic string, messages []string) {
	f.T.Helper()
	
	for _, msg := range messages {
		err := f.Producer.Publish(f.Ctx, topic, []byte(msg), nil)
		require.NoError(f.T, err, "Should publish message %s", msg)
	}
}

// AssertMessagesReceived waits for and validates received messages
func (f *BrokerFixture) AssertMessagesReceived(msgChan <-chan *queue.Message, expectedCount int, timeout time.Duration) {
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

// CreateLogger creates a test logger with the given prefix
func CreateLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, "["+prefix+"] ", log.LstdFlags)
}

// Constants for common test values
const (
	ConsumerStartupDelay = 200 * time.Millisecond
	DefaultTestTimeout   = 5 * time.Second
)
