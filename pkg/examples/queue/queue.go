// Package queue provides interfaces for message queuing operations
package queue

import (
	"context"
	"time"
)

// Message represents a message in the queue
type Message struct {
	ID        string            `json:"id"`
	Topic     string            `json:"topic"`
	Payload   []byte            `json:"payload"`
	Headers   map[string]string `json:"headers"`
	Timestamp time.Time         `json:"timestamp"`
}

// Queue interface defines the basic queue operations
type Queue interface {
	// Enqueue adds a message to the specified topic
	Enqueue(ctx context.Context, topic string, message *Message) error
	
	// Dequeue retrieves a message from the specified topic
	// Returns nil if no message is available
	Dequeue(ctx context.Context, topic string) (*Message, error)
	
	// Size returns the number of messages in the specified topic
	Size(ctx context.Context, topic string) (int, error)
	
	// Topics returns all available topics
	Topics(ctx context.Context) ([]string, error)
	
	// Close closes the queue and releases resources
	Close() error
}

// Producer interface defines message publishing operations
type Producer interface {
	// Publish sends a message to the specified topic
	Publish(ctx context.Context, topic string, payload []byte, headers map[string]string) error
	
	// Close closes the producer and releases resources
	Close() error
}

// MessageHandler is a function type for handling received messages
type MessageHandler func(ctx context.Context, message *Message) error

// Consumer interface defines message consumption operations
type Consumer interface {
	// Subscribe starts consuming messages from the specified topic
	// The handler function will be called for each received message
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	
	// Unsubscribe stops consuming messages from the specified topic
	Unsubscribe(ctx context.Context, topic string) error
	
	// Close closes the consumer and releases resources
	Close() error
}

