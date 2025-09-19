# Queue System

A comprehensive in-memory queue system with producer-consumer patterns implemented in Go.

## Architecture

The queue system consists of four main packages:

### 1. `queue` - Core Interfaces
Contains the fundamental interfaces:
- `Queue`: Basic queue operations (Enqueue, Dequeue, Size, Topics, Close)
- `Producer`: Message publishing interface
- `Consumer`: Message consumption interface with subscription support
- `Message`: Standardized message structure with ID, Topic, Payload, Headers, and Timestamp

### 2. `inmemory` - In-Memory Queue Implementation
Implements the `Queue` interface using:
- Thread-safe in-memory storage with mutexes
- Buffered channels for each topic (capacity: 1000 messages)
- Graceful shutdown handling

### 3. `broker` - Producer and Consumer Implementations
- `QueueProducer`: Implements `Producer` interface using any `Queue` implementation
- `QueueConsumer`: Implements `Consumer` interface with subscription management and polling

### 4. `example` - Working Example Services
- `ProducerService`: Generates order messages every 2 seconds
- `ConsumerService`: Processes order messages with business logic
- `RunExample()`: Demonstrates the complete system working together



