package queue

import (
	"context"
	"errors"
	"sync"
)

// Mock is a simple in-memory queue implementation for testing
// It implements the Queue interface and can be used by any package for testing
type Mock struct {
	topics map[string]chan *Message
	closed bool
	mutex  sync.RWMutex
}

// NewMock creates a new mock queue for testing
func NewMock() *Mock {
	return &Mock{
		topics: make(map[string]chan *Message),
		closed: false,
	}
}

// Enqueue adds a message to the specified topic
func (q *Mock) Enqueue(ctx context.Context, topic string, message *Message) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.closed {
		return errors.New("queue is closed")
	}

	if _, exists := q.topics[topic]; !exists {
		q.topics[topic] = make(chan *Message, 100)
	}

	select {
	case q.topics[topic] <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Dequeue retrieves a message from the specified topic
func (q *Mock) Dequeue(ctx context.Context, topic string) (*Message, error) {
	q.mutex.RLock()
	ch, exists := q.topics[topic]
	closed := q.closed
	q.mutex.RUnlock()

	if closed {
		return nil, errors.New("queue is closed")
	}

	if !exists {
		return nil, nil
	}

	select {
	case msg := <-ch:
		return msg, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, nil
	}
}

// Size returns the number of messages in the specified topic
func (q *Mock) Size(ctx context.Context, topic string) (int, error) {
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

// Topics returns a list of all topics that have been created
func (q *Mock) Topics(ctx context.Context) ([]string, error) {
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

// Close closes the mock queue and all its topic channels
func (q *Mock) Close() error {
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
