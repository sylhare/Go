package inmemory

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syl/Go/pkg/examples/queue"
	"github.com/syl/Go/pkg/examples/queue/inmemory/testutils"
)

func assertQueueOperationErrors(t *testing.T, fixture *testutils.BaseFixture, topic string, msg *queue.Message) {
	t.Helper()
	
	err := fixture.Queue.Enqueue(fixture.Ctx, topic, msg)
	assert.Error(t, err, "Should error when enqueueing to closed queue")

	_, err = fixture.Queue.Dequeue(fixture.Ctx, topic)
	assert.Error(t, err, "Should error when dequeueing from closed queue")

	_, err = fixture.Queue.Size(fixture.Ctx, topic)
	assert.Error(t, err, "Should error when getting size of closed queue")

	_, err = fixture.Queue.Topics(fixture.Ctx)
	assert.Error(t, err, "Should error when getting topics of closed queue")
}

func TestInMemoryQueue(t *testing.T) {
	t.Run("EnqueueDequeue", func(t *testing.T) {
		q := NewInMemoryQueue()
		fixture := testutils.NewBaseFixture(t, q)
		topic := "test-topic"

		msg := fixture.CreateMessage("test-id", topic, []byte("test payload"))

		err := fixture.Queue.Enqueue(fixture.Ctx, topic, msg)
		require.NoError(t, err, "Should enqueue message successfully")

		fixture.AssertQueueSize(topic, 1, "Queue should contain exactly one message")

		dequeuedMsg, err := fixture.Queue.Dequeue(fixture.Ctx, topic)
		require.NoError(t, err, "Should dequeue message")
		require.NotNil(t, dequeuedMsg, "Message should not be nil")

		assert.Equal(t, msg.ID, dequeuedMsg.ID, "Message ID should match")
		assert.Equal(t, msg.Payload, dequeuedMsg.Payload, "Payload should match")
		assert.Equal(t, msg.Headers, dequeuedMsg.Headers, "Headers should match")
	})

	t.Run("DequeueEmpty", func(t *testing.T) {
		q := NewInMemoryQueue()
		fixture := testutils.NewBaseFixture(t, q)
		topic := "empty-topic"

		msg, err := fixture.Queue.Dequeue(fixture.Ctx, topic)
		require.NoError(t, err, "Should not error when dequeuing from empty queue")
		assert.Nil(t, msg, "Should return nil message for empty queue")
	})

	t.Run("Topics", func(t *testing.T) {
		q := NewInMemoryQueue()
		fixture := testutils.NewBaseFixture(t, q)

		t.Run("InitiallyEmpty", func(t *testing.T) {
			topics, err := fixture.Queue.Topics(fixture.Ctx)
			require.NoError(t, err, "Should get topics successfully")
			assert.Empty(t, topics, "Should have no topics initially")
		})

		t.Run("AfterEnqueue", func(t *testing.T) {
			msg1 := fixture.CreateMessage("1", "topic1", []byte("payload1"))
			msg2 := fixture.CreateMessage("2", "topic2", []byte("payload2"))

			err := fixture.Queue.Enqueue(fixture.Ctx, "topic1", msg1)
			require.NoError(t, err, "Should enqueue to topic1")

			err = fixture.Queue.Enqueue(fixture.Ctx, "topic2", msg2)
			require.NoError(t, err, "Should enqueue to topic2")

			topics, err := fixture.Queue.Topics(fixture.Ctx)
			require.NoError(t, err, "Should get topics successfully")
			assert.Len(t, topics, 2, "Should have exactly 2 topics")

			fixture.AssertTopicsContain("topic1", "topic2")
		})
	})

	t.Run("Close", func(t *testing.T) {
		q := NewInMemoryQueue()
		fixture := testutils.NewBaseFixture(t, q)
		topic := "test-topic"

		msg := fixture.CreateMessage("test", topic, []byte("test"))
		err := fixture.Queue.Enqueue(fixture.Ctx, topic, msg)
		require.NoError(t, err, "Should enqueue message before closing")

		err = fixture.Queue.Close()
		require.NoError(t, err, "Should close queue successfully")

		t.Run("OperationsAfterClose", func(t *testing.T) {
			assertQueueOperationErrors(t, fixture, topic, msg)
		})
	})

	t.Run("Concurrent", func(t *testing.T) {
		q := NewInMemoryQueue()
		fixture := testutils.NewBaseFixture(t, q)
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
					fixture.Queue.Enqueue(fixture.Ctx, topic, msg)
				}
				done <- struct{}{}
			}(i)
		}

		for i := 0; i < numProducers; i++ {
			<-done
		}

		expectedSize := numProducers * numMessages
		fixture.AssertQueueSize(topic, expectedSize, "Should have correct number of messages")

		messagesReceived := 0
		for {
			msg, err := fixture.Queue.Dequeue(fixture.Ctx, topic)
			require.NoError(t, err, "Should dequeue messages without error")
			if msg == nil {
				break
			}
			messagesReceived++
		}

		assert.Equal(t, expectedSize, messagesReceived, "Should receive all messages")
	})
}
