package inmemory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/syl/Go/pkg/examples/queue"
)

func TestInMemoryQueue_EnqueueDequeue(t *testing.T) {
	q := NewInMemoryQueue()
	defer q.Close()

	ctx := context.Background()
	topic := "test-topic"

	msg := &queue.Message{
		ID:        "test-id",
		Topic:     topic,
		Payload:   []byte("test payload"),
		Headers:   map[string]string{"key": "value"},
		Timestamp: time.Now(),
	}

	err := q.Enqueue(ctx, topic, msg)
	if err != nil {
		t.Fatalf("Failed to enqueue message: %v", err)
	}

	size, err := q.Size(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to get size: %v", err)
	}
	if size != 1 {
		t.Fatalf("Expected size 1, got %d", size)
	}

	dequeuedMsg, err := q.Dequeue(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to dequeue message: %v", err)
	}

	if dequeuedMsg == nil {
		t.Fatal("Expected message, got nil")
	}

	if dequeuedMsg.ID != msg.ID {
		t.Fatalf("Expected message ID %s, got %s", msg.ID, dequeuedMsg.ID)
	}

	if string(dequeuedMsg.Payload) != string(msg.Payload) {
		t.Fatalf("Expected payload %s, got %s", string(msg.Payload), string(dequeuedMsg.Payload))
	}
}

func TestInMemoryQueue_DequeueEmpty(t *testing.T) {
	q := NewInMemoryQueue()
	defer q.Close()

	ctx := context.Background()
	topic := "empty-topic"

	msg, err := q.Dequeue(ctx, topic)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if msg != nil {
		t.Fatal("Expected nil message from empty queue")
	}
}

func TestInMemoryQueue_Topics(t *testing.T) {
	q := NewInMemoryQueue()
	defer q.Close()

	ctx := context.Background()

	topics, err := q.Topics(ctx)
	if err != nil {
		t.Fatalf("Failed to get topics: %v", err)
	}
	if len(topics) != 0 {
		t.Fatalf("Expected 0 topics, got %d", len(topics))
	}

	msg1 := &queue.Message{ID: "1", Topic: "topic1", Payload: []byte("payload1"), Timestamp: time.Now()}
	msg2 := &queue.Message{ID: "2", Topic: "topic2", Payload: []byte("payload2"), Timestamp: time.Now()}

	err = q.Enqueue(ctx, "topic1", msg1)
	if err != nil {
		t.Fatalf("Failed to enqueue to topic1: %v", err)
	}

	err = q.Enqueue(ctx, "topic2", msg2)
	if err != nil {
		t.Fatalf("Failed to enqueue to topic2: %v", err)
	}

	topics, err = q.Topics(ctx)
	if err != nil {
		t.Fatalf("Failed to get topics: %v", err)
	}

	if len(topics) != 2 {
		t.Fatalf("Expected 2 topics, got %d", len(topics))
	}

	topicMap := make(map[string]bool)
	for _, topic := range topics {
		topicMap[topic] = true
	}

	if !topicMap["topic1"] || !topicMap["topic2"] {
		t.Fatal("Expected both topic1 and topic2 to be present")
	}
}

func TestInMemoryQueue_Close(t *testing.T) {
	q := NewInMemoryQueue()

	ctx := context.Background()
	topic := "test-topic"

	msg := &queue.Message{ID: "test", Topic: topic, Payload: []byte("test"), Timestamp: time.Now()}
	err := q.Enqueue(ctx, topic, msg)
	if err != nil {
		t.Fatalf("Failed to enqueue message: %v", err)
	}

	err = q.Close()
	if err != nil {
		t.Fatalf("Failed to close queue: %v", err)
	}

	err = q.Enqueue(ctx, topic, msg)
	if err == nil {
		t.Fatal("Expected error when enqueueing to closed queue")
	}

	_, err = q.Dequeue(ctx, topic)
	if err == nil {
		t.Fatal("Expected error when dequeueing from closed queue")
	}

	_, err = q.Size(ctx, topic)
	if err == nil {
		t.Fatal("Expected error when getting size of closed queue")
	}

	_, err = q.Topics(ctx)
	if err == nil {
		t.Fatal("Expected error when getting topics of closed queue")
	}
}

func TestInMemoryQueue_Concurrent(t *testing.T) {
	q := NewInMemoryQueue()
	defer q.Close()

	ctx := context.Background()
	topic := "concurrent-topic"
	numProducers := 5
	numMessages := 10

	done := make(chan struct{})
	for i := 0; i < numProducers; i++ {
		go func(producerID int) {
			for j := 0; j < numMessages; j++ {
				msg := &queue.Message{
					ID:        fmt.Sprintf("producer-%d-msg-%d", producerID, j),
					Topic:     topic,
					Payload:   []byte(fmt.Sprintf("payload from producer %d, message %d", producerID, j)),
					Timestamp: time.Now(),
				}
				q.Enqueue(ctx, topic, msg)
			}
			done <- struct{}{}
		}(i)
	}

	for i := 0; i < numProducers; i++ {
		<-done
	}

	size, err := q.Size(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to get size: %v", err)
	}

	expectedSize := numProducers * numMessages
	if size != expectedSize {
		t.Fatalf("Expected size %d, got %d", expectedSize, size)
	}

	messagesReceived := 0
	for {
		msg, err := q.Dequeue(ctx, topic)
		if err != nil {
			t.Fatalf("Failed to dequeue message: %v", err)
		}
		if msg == nil {
			break
		}
		messagesReceived++
	}

	if messagesReceived != expectedSize {
		t.Fatalf("Expected to receive %d messages, got %d", expectedSize, messagesReceived)
	}
}
