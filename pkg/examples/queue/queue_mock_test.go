package queue

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMock(t *testing.T) {
	t.Run("EnqueueDequeue", func(t *testing.T) {
		q := NewMock()
		defer q.Close()

		ctx := context.Background()
		topic := "test-topic"

		msg := &Message{
			ID:        "test-1",
			Topic:     topic,
			Payload:   []byte("test message"),
			Headers:   map[string]string{"key": "value"},
			Timestamp: time.Now(),
		}

		t.Run("EnqueueMessage", func(t *testing.T) {
			err := q.Enqueue(ctx, topic, msg)
			require.NoError(t, err, "Should enqueue message successfully")

			size, err := q.Size(ctx, topic)
			require.NoError(t, err, "Should get queue size")
			assert.Equal(t, 1, size, "Queue should have one message")
		})

		t.Run("DequeueMessage", func(t *testing.T) {
			dequeued, err := q.Dequeue(ctx, topic)
			require.NoError(t, err, "Should dequeue message successfully")
			require.NotNil(t, dequeued, "Dequeued message should not be nil")

			assert.Equal(t, msg.ID, dequeued.ID, "Message ID should match")
			assert.Equal(t, msg.Topic, dequeued.Topic, "Message topic should match")
			assert.Equal(t, msg.Payload, dequeued.Payload, "Message payload should match")
			assert.Equal(t, msg.Headers, dequeued.Headers, "Message headers should match")

			size, err := q.Size(ctx, topic)
			require.NoError(t, err, "Should get queue size")
			assert.Equal(t, 0, size, "Queue should be empty after dequeue")
		})
	})

	t.Run("DequeueEmpty", func(t *testing.T) {
		q := NewMock()
		defer q.Close()

		ctx := context.Background()
		topic := "empty-topic"

		t.Run("DequeueFromEmptyTopic", func(t *testing.T) {
			msg, err := q.Dequeue(ctx, topic)
			require.NoError(t, err, "Should not error when dequeuing from empty topic")
			assert.Nil(t, msg, "Should return nil for empty topic")
		})

		t.Run("DequeueFromNonExistentTopic", func(t *testing.T) {
			msg, err := q.Dequeue(ctx, "non-existent")
			require.NoError(t, err, "Should not error for non-existent topic")
			assert.Nil(t, msg, "Should return nil for non-existent topic")
		})
	})

	t.Run("Size", func(t *testing.T) {
		q := NewMock()
		defer q.Close()

		ctx := context.Background()
		topic := "size-test"

		t.Run("EmptyTopicSize", func(t *testing.T) {
			size, err := q.Size(ctx, topic)
			require.NoError(t, err, "Should get size for empty topic")
			assert.Equal(t, 0, size, "Empty topic should have size 0")
		})

		t.Run("MultipleMessages", func(t *testing.T) {
			for i := 0; i < 5; i++ {
				msg := &Message{
					ID:      "msg-" + string(rune('1'+i)),
					Topic:   topic,
					Payload: []byte("test"),
					Headers: map[string]string{},
				}
				err := q.Enqueue(ctx, topic, msg)
				require.NoError(t, err, "Should enqueue message %d", i+1)
			}

			size, err := q.Size(ctx, topic)
			require.NoError(t, err, "Should get queue size")
			assert.Equal(t, 5, size, "Should have 5 messages in queue")
		})
	})

	t.Run("Topics", func(t *testing.T) {
		q := NewMock()
		defer q.Close()

		ctx := context.Background()

		t.Run("InitiallyEmpty", func(t *testing.T) {
			topics, err := q.Topics(ctx)
			require.NoError(t, err, "Should get topics list")
			assert.Len(t, topics, 0, "Should have no topics initially")
		})

		t.Run("AfterEnqueue", func(t *testing.T) {
			testTopics := []string{"topic1", "topic2", "topic3"}
			
			for _, topic := range testTopics {
				msg := &Message{
					ID:      "msg-" + topic,
					Topic:   topic,
					Payload: []byte("test"),
					Headers: map[string]string{},
				}
				err := q.Enqueue(ctx, topic, msg)
				require.NoError(t, err, "Should enqueue to %s", topic)
			}

			topics, err := q.Topics(ctx)
			require.NoError(t, err, "Should get topics list")
			assert.Len(t, topics, len(testTopics), "Should have correct number of topics")

			topicMap := make(map[string]bool)
			for _, topic := range topics {
				topicMap[topic] = true
			}

			for _, expectedTopic := range testTopics {
				assert.True(t, topicMap[expectedTopic], "Should contain topic %s", expectedTopic)
			}
		})
	})

	t.Run("Close", func(t *testing.T) {
		q := NewMock()
		ctx := context.Background()

		msg := &Message{
			ID:      "test-close",
			Topic:   "close-topic",
			Payload: []byte("test"),
			Headers: map[string]string{},
		}
		err := q.Enqueue(ctx, "close-topic", msg)
		require.NoError(t, err, "Should enqueue before close")

		t.Run("CloseSuccessfully", func(t *testing.T) {
			err := q.Close()
			require.NoError(t, err, "Should close successfully")
		})

		t.Run("OperationsAfterClose", func(t *testing.T) {
			err := q.Enqueue(ctx, "test", msg)
			assert.Error(t, err, "Enqueue should fail after close")

			_, err = q.Dequeue(ctx, "test")
			assert.Error(t, err, "Dequeue should fail after close")

			_, err = q.Size(ctx, "test")
			assert.Error(t, err, "Size should fail after close")

			_, err = q.Topics(ctx)
			assert.Error(t, err, "Topics should fail after close")

			err = q.Close()
			require.NoError(t, err, "Multiple closes should be safe")
		})
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		q := NewMock()
		defer q.Close()

		t.Run("EnqueueWithCancelledContext", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			msg := &Message{ID: "test", Topic: "test", Payload: []byte("test")}
			err := q.Enqueue(ctx, "test", msg)
			assert.Error(t, err, "Should return context error")
			assert.Equal(t, context.Canceled, err, "Should return context.Canceled")
		})

		t.Run("DequeueWithCancelledContext", func(t *testing.T) {
			msg := &Message{ID: "test", Topic: "test", Payload: []byte("test")}
			err := q.Enqueue(context.Background(), "test", msg)
			require.NoError(t, err, "Should enqueue successfully")

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err = q.Dequeue(ctx, "test")
		})
	})

	t.Run("Interface", func(t *testing.T) {
		t.Run("ImplementsQueueInterface", func(t *testing.T) {
			var _ Queue = (*Mock)(nil)
			
			q := NewMock()
			defer q.Close()

			var queue Queue = q
			assert.NotNil(t, queue, "Should implement Queue interface")
		})
	})
}
