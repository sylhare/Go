// Package inmemory provides an in-memory implementation of the queue interface
package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/syl/Go/pkg/examples/queue"
)

// InMemoryQueue implements the Queue interface using in-memory storage
type InMemoryQueue struct {
	mu     sync.RWMutex
	topics map[string]chan *queue.Message
	closed bool
}

// NewInMemoryQueue creates a new in-memory queue
func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		topics: make(map[string]chan *queue.Message),
		closed: false,
	}
}

// Enqueue adds a message to the specified topic
func (q *InMemoryQueue) Enqueue(ctx context.Context, topic string, message *queue.Message) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return fmt.Errorf("queue is closed")
	}

	if _, exists := q.topics[topic]; !exists {
		q.topics[topic] = make(chan *queue.Message, 1000)
	}

	select {
	case q.topics[topic] <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("topic %s queue is full", topic)
	}
}

// Dequeue retrieves a message from the specified topic
func (q *InMemoryQueue) Dequeue(ctx context.Context, topic string) (*queue.Message, error) {
	q.mu.RLock()
	topicChan, exists := q.topics[topic]
	closed := q.closed
	q.mu.RUnlock()

	if closed {
		return nil, fmt.Errorf("queue is closed")
	}

	if !exists {
		return nil, nil
	}

	select {
	case message := <-topicChan:
		return message, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, nil
	}
}

// Size returns the number of messages in the specified topic
func (q *InMemoryQueue) Size(ctx context.Context, topic string) (int, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.closed {
		return 0, fmt.Errorf("queue is closed")
	}

	if topicChan, exists := q.topics[topic]; exists {
		return len(topicChan), nil
	}

	return 0, nil
}

// Topics returns all available topics
func (q *InMemoryQueue) Topics(ctx context.Context) ([]string, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.closed {
		return nil, fmt.Errorf("queue is closed")
	}

	topics := make([]string, 0, len(q.topics))
	for topic := range q.topics {
		topics = append(topics, topic)
	}

	return topics, nil
}

// Close closes the queue and releases resources
func (q *InMemoryQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil
	}

	q.closed = true

	for _, ch := range q.topics {
		close(ch)
	}

	q.topics = make(map[string]chan *queue.Message)

	return nil
}
