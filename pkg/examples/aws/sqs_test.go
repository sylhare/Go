package localstack

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
)

func TestSQS(t *testing.T) {
	ctx := context.Background()

	localStack, err := New()
	if err != nil {
		t.Fatalf("failed to start LocalStack: %s", err)
	}
	defer localStack.Terminate()

	testSQS := sqs.NewFromConfig(localStack.Config)
	var queueUrl string

	t.Run("CreateQueue", func(t *testing.T) {
		queueName := "test-queue"
		output, err := testSQS.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: &queueName,
		})
		if err != nil {
			t.Fatalf("failed to create queue: %s", err)
		}
		queueUrl = *output.QueueUrl
		assert.NotEmpty(t, queueUrl, "queue URL should not be empty")
	})

	t.Run("SendMessage", func(t *testing.T) {
		messageBody := "Hello, SQS!"
		_, err := testSQS.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    &queueUrl,
			MessageBody: &messageBody,
		})
		if err != nil {
			t.Fatalf("failed to send message: %s", err)
		}
	})

	t.Run("ReceiveMessage", func(t *testing.T) {
		output, err := testSQS.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &queueUrl,
			MaxNumberOfMessages: 10,
		})
		if err != nil {
			t.Fatalf("failed to receive message: %s", err)
		}
		messages := make([]string, len(output.Messages))
		for i, msg := range output.Messages {
			messages[i] = *msg.Body
		}
		assert.NotEmpty(t, messages, "should receive at least one message")
		assert.Equal(t, "Hello, SQS!", messages[0], "message body should match")
	})
}
