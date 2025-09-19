package broker

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/syl/Go/pkg/examples/queue"
	"github.com/syl/Go/pkg/examples/queue/inmemory"
)

func TestQueueProducer_Publish(t *testing.T) {
	q := inmemory.NewInMemoryQueue()
	defer q.Close()

	producer := NewQueueProducer(q)
	defer producer.Close()

	ctx := context.Background()
	topic := "test-topic"
	payload := []byte("test message")
	headers := map[string]string{"key": "value"}

	err := producer.Publish(ctx, topic, payload, headers)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	size, err := q.Size(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to get queue size: %v", err)
	}

	if size != 1 {
		t.Fatalf("Expected queue size 1, got %d", size)
	}

	msg, err := q.Dequeue(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to dequeue message: %v", err)
	}

	if msg == nil {
		t.Fatal("Expected message, got nil")
	}

	if string(msg.Payload) != string(payload) {
		t.Fatalf("Expected payload %s, got %s", string(payload), string(msg.Payload))
	}

	if msg.Headers["key"] != "value" {
		t.Fatal("Expected header not found or incorrect")
	}

	if msg.Topic != topic {
		t.Fatalf("Expected topic %s, got %s", topic, msg.Topic)
	}
}

func TestQueueConsumer_Subscribe(t *testing.T) {
	q := inmemory.NewInMemoryQueue()
	defer q.Close()

	producer := NewQueueProducer(q)
	defer producer.Close()

	consumer := NewQueueConsumer(q)
	defer consumer.Close()

	ctx := context.Background()
	topic := "test-topic"

	receivedMessages := make(chan *queue.Message, 10)

	handler := func(ctx context.Context, message *queue.Message) error {
		receivedMessages <- message
		return nil
	}

	err := consumer.Subscribe(ctx, topic, handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	testMessages := []string{"message1", "message2", "message3"}
	for _, msg := range testMessages {
		err := producer.Publish(ctx, topic, []byte(msg), nil)
		if err != nil {
			t.Fatalf("Failed to publish message %s: %v", msg, err)
		}
	}

	receivedCount := 0
	timeout := time.After(5 * time.Second)

	for receivedCount < len(testMessages) {
		select {
		case msg := <-receivedMessages:
			t.Logf("Received message: %s", string(msg.Payload))
			receivedCount++
		case <-timeout:
			t.Fatalf("Timeout waiting for messages, received %d of %d", receivedCount, len(testMessages))
		}
	}

	err = consumer.Unsubscribe(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}
}

func TestQueueConsumer_MultipleSubscriptions(t *testing.T) {
	q := inmemory.NewInMemoryQueue()
	defer q.Close()

	producer := NewQueueProducer(q)
	defer producer.Close()

	consumer := NewQueueConsumer(q)
	defer consumer.Close()

	ctx := context.Background()
	topic1 := "topic1"
	topic2 := "topic2"

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

	err := consumer.Subscribe(ctx, topic1, handler1)
	if err != nil {
		t.Fatalf("Failed to subscribe to topic1: %v", err)
	}

	err = consumer.Subscribe(ctx, topic2, handler2)
	if err != nil {
		t.Fatalf("Failed to subscribe to topic2: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	err = producer.Publish(ctx, topic1, []byte("message for topic1"), nil)
	if err != nil {
		t.Fatalf("Failed to publish to topic1: %v", err)
	}

	err = producer.Publish(ctx, topic2, []byte("message for topic2"), nil)
	if err != nil {
		t.Fatalf("Failed to publish to topic2: %v", err)
	}

	timeout := time.After(5 * time.Second)
	receivedTopic1 := false
	receivedTopic2 := false

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
}

func TestQueueConsumer_ConcurrentConsumers(t *testing.T) {
	q := inmemory.NewInMemoryQueue()
	defer q.Close()

	producer := NewQueueProducer(q)
	defer producer.Close()

	ctx := context.Background()
	topic := "concurrent-topic"
	numConsumers := 3
	numMessages := 10

	allMessages := make(chan *queue.Message, numMessages)
	var wg sync.WaitGroup

	consumers := make([]*QueueConsumer, numConsumers)
	for i := 0; i < numConsumers; i++ {
		consumers[i] = NewQueueConsumer(q)
		
		handler := func(ctx context.Context, message *queue.Message) error {
			allMessages <- message
			return nil
		}

		err := consumers[i].Subscribe(ctx, topic, handler)
		if err != nil {
			t.Fatalf("Failed to subscribe consumer %d: %v", i, err)
		}
	}

	time.Sleep(200 * time.Millisecond)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numMessages; i++ {
			err := producer.Publish(ctx, topic, []byte(fmt.Sprintf("message-%d", i)), nil)
			if err != nil {
				t.Errorf("Failed to publish message %d: %v", i, err)
			}
		}
	}()

	wg.Wait()

	receivedCount := 0
	timeout := time.After(10 * time.Second)

	for receivedCount < numMessages {
		select {
		case msg := <-allMessages:
			t.Logf("Received message: %s", string(msg.Payload))
			receivedCount++
		case <-timeout:
			t.Fatalf("Timeout waiting for messages, received %d of %d", receivedCount, numMessages)
		}
	}

	for i := range consumers {
		consumers[i].Close()
	}
}
