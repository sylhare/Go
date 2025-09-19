package testutils

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syl/Go/pkg/examples/queue"
)

// ServiceFixture extends BrokerFixture with service-specific functionality
type ServiceFixture struct {
	*BrokerFixture
}

// NewServiceFixture creates a fixture for service testing
func NewServiceFixture(t *testing.T, queue queue.Queue, producer queue.Producer, consumer queue.Consumer) *ServiceFixture {
	t.Helper()
	
	brokerFixture := NewBrokerFixture(t, queue, producer, consumer)
	
	return &ServiceFixture{
		BrokerFixture: brokerFixture,
	}
}

// CreateTestMessageWithPayload creates a test message with custom payload
func (f *ServiceFixture) CreateTestMessageWithPayload(id, topic string, payload []byte, headers map[string]string) *queue.Message {
	f.T.Helper()
	
	return &queue.Message{
		ID:        "test-msg-" + id,
		Topic:     topic,
		Payload:   payload,
		Headers:   headers,
		Timestamp: time.Now(),
	}
}

// RunServiceWithTimeout runs a service function with a timeout
func (f *ServiceFixture) RunServiceWithTimeout(serviceFunc func(context.Context) error, timeout time.Duration) {
	f.T.Helper()
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := serviceFunc(ctx)
		assert.True(f.T, err == nil || err == context.DeadlineExceeded, "Service should complete or timeout gracefully")
	}()

	wg.Wait()
}

// AssertMessageHeaders validates that a message has expected headers
func (f *ServiceFixture) AssertMessageHeaders(msg *queue.Message, expectedHeaders map[string]string) {
	f.T.Helper()
	
	for key, expectedValue := range expectedHeaders {
		actualValue, exists := msg.Headers[key]
		require.True(f.T, exists, "Header %s should exist", key)
		assert.Equal(f.T, expectedValue, actualValue, "Header %s should match expected value", key)
	}
}

// AssertValidJSON validates that message payload is valid JSON
func (f *ServiceFixture) AssertValidJSON(msg *queue.Message) {
	f.T.Helper()
	
	var data interface{}
	err := json.Unmarshal(msg.Payload, &data)
	require.NoError(f.T, err, "Payload should be valid JSON")
}
