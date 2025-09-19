package broker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/syl/Go/pkg/examples/queue"
)

// QueueConsumer implements the Consumer interface using a Queue
type QueueConsumer struct {
	queue         queue.Queue
	subscriptions map[string]*subscription
	mu            sync.RWMutex
	closed        bool
}

// subscription represents an active subscription to a topic
type subscription struct {
	topic     string
	handler   queue.MessageHandler
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// NewQueueConsumer creates a new consumer that uses the provided queue
func NewQueueConsumer(q queue.Queue) *QueueConsumer {
	return &QueueConsumer{
		queue:         q,
		subscriptions: make(map[string]*subscription),
		closed:        false,
	}
}

// Subscribe starts consuming messages from the specified topic
func (c *QueueConsumer) Subscribe(ctx context.Context, topic string, handler queue.MessageHandler) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return fmt.Errorf("consumer is closed")
	}

	if _, exists := c.subscriptions[topic]; exists {
		return fmt.Errorf("already subscribed to topic: %s", topic)
	}

	subCtx, cancel := context.WithCancel(ctx)
	
	sub := &subscription{
		topic:   topic,
		handler: handler,
		ctx:     subCtx,
		cancel:  cancel,
	}

	c.subscriptions[topic] = sub

	sub.wg.Add(1)
	go c.consumeMessages(sub)

	return nil
}

// Unsubscribe stops consuming messages from the specified topic
func (c *QueueConsumer) Unsubscribe(ctx context.Context, topic string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	sub, exists := c.subscriptions[topic]
	if !exists {
		return fmt.Errorf("not subscribed to topic: %s", topic)
	}

	sub.cancel()
	
	sub.wg.Wait()
	
	delete(c.subscriptions, topic)

	return nil
}

// Close closes the consumer and releases resources
func (c *QueueConsumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	for _, sub := range c.subscriptions {
		sub.cancel()
	}

	for _, sub := range c.subscriptions {
		sub.wg.Wait()
	}

	c.subscriptions = make(map[string]*subscription)

	return nil
}

// consumeMessages continuously polls for messages from the queue
func (c *QueueConsumer) consumeMessages(sub *subscription) {
	defer sub.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sub.ctx.Done():
			return
		case <-ticker.C:
			message, err := c.queue.Dequeue(sub.ctx, sub.topic)
			if err != nil {
				continue
			}

			if message != nil {
				if err := sub.handler(sub.ctx, message); err != nil {
					continue
				}
			}
		}
	}
}
