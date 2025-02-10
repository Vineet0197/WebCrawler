package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// KafkaStorage represents the Kafka storage structure
type KafkaStorage struct {
	Client  *kgo.Client
	Brokers []string
	Topic   string
	GroupID string
}

// NewKafkaStorage creates a new Kafka storage
func NewKafkaStorage(brokers []string, topic string, groupID string, username, password string) (*KafkaStorage, error) {
	// Set up SASL authentication using SCRAM-SHA-512
	scramAuth := scram.Auth{
		User: username,
		Pass: password,
	}
	scramClient := scramAuth.AsSha512Mechanism() // Redpanda Cloud requires SCRAM-SHA-512

	// Configure Kafka client with both producer and consumer
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.SASL(scramClient),
		kgo.DialTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}), // Enable TLS for SASL_SSL
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup(groupID),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client: %v", err)
	}

	return &KafkaStorage{Client: client, Topic: topic}, nil
}

// Produce sends a message to the Kafka topic
func (k *KafkaStorage) Produce(message string) error {
	// Send the message to the topic
	record := &kgo.Record{Topic: k.Topic, Value: []byte(message)}
	ctx := context.Background()
	err := k.Client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	log.Printf("Message sent to topic %s", k.Topic)
	return nil
}

// Consume receives a message from the Kafka topic
func (k *KafkaStorage) Consume(ctx context.Context, topic string, handler func(message string)) {
	// Consume messages from the topic
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Start the consumer group in a separate go routine
	go func() {
		defer wg.Done()
		for {
			fetches := k.Client.PollFetches(ctx)
			if fetches.IsClientClosed() {
				log.Println("Kafka consumer closed, stopping consumption")
				time.Sleep(1 * time.Second)
			}

			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()
				handler(string(record.Value)) // Pass the message to the handler function
			}

			// Check if the context is done
			if ctx.Err() != nil {
				return
			}
		}
	}()

	// Wait for the consumer group to finish
	wg.Wait()
}

// Close closes the Kafka producer and consumer
func (k *KafkaStorage) Close() {
	// Close the Kafka Client
	k.Client.Close()
	log.Printf("Kafka client closed successfully")
}
