package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/syl/Go/pkg/examples/queue/broker"
	"github.com/syl/Go/pkg/examples/queue/example"
	"github.com/syl/Go/pkg/examples/queue/inmemory"
)

func main() {
	producerLogger := log.New(os.Stdout, "[PRODUCER] ", log.LstdFlags|log.Lshortfile)
	consumerLogger := log.New(os.Stdout, "[CONSUMER] ", log.LstdFlags|log.Lshortfile)
	mainLogger := log.New(os.Stdout, "[MAIN] ", log.LstdFlags|log.Lshortfile)

	mainLogger.Println("Starting queue system example...")

	queue := inmemory.NewInMemoryQueue()
	defer queue.Close()

	producer := broker.NewQueueProducer(queue)
	consumer := broker.NewQueueConsumer(queue)

	producerService := example.NewProducerService(producer, producerLogger)
	consumerService := example.NewConsumerService(consumer, consumerLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := consumerService.Start(ctx); err != nil && err != context.Canceled {
			mainLogger.Printf("Consumer service error: %v", err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := producerService.Start(ctx); err != nil && err != context.Canceled {
			mainLogger.Printf("Producer service error: %v", err)
		}
	}()

	mainLogger.Println("Services started. Press Ctrl+C to stop...")

	<-sigChan
	mainLogger.Println("Shutting down services...")

	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		mainLogger.Println("All services shut down gracefully")
	case <-time.After(10 * time.Second):
		mainLogger.Println("Shutdown timeout reached")
	}

	if err := producerService.Stop(); err != nil {
		mainLogger.Printf("Error stopping producer service: %v", err)
	}

	if err := consumerService.Stop(); err != nil {
		mainLogger.Printf("Error stopping consumer service: %v", err)
	}

	mainLogger.Println("Example completed")
}
