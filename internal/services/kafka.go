package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type consumerGroupHandler struct {
	handler func(message string)
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.handler(string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}

// KafkaStorage represents the Kafka storage structure
type KafkaStorage struct {
	Producer sarama.SyncProducer
	Consumer sarama.ConsumerGroup
	Brokers  []string
	Topic    string
	GroupID  string
}

// NewKafkaStorage creates a new Kafka storage
func NewKafkaStorage(brokers []string, topic string, groupID string) (*KafkaStorage, error) {
	// Configure the Kafka producer
	producerCfg := sarama.NewConfig()
	producerCfg.Producer.RequiredAcks = sarama.WaitForAll
	producerCfg.Producer.Retry.Max = 5
	producerCfg.Producer.Return.Successes = true

	// Creates a new producer
	producer, err := sarama.NewSyncProducer(brokers, producerCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	// Configure the Kafka consumer
	consumerCfg := sarama.NewConfig()
	consumerCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerCfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()

	// Creates a new consumer group
	consumer, err := sarama.NewConsumerGroup(brokers, groupID, consumerCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer group: %v", err)
	}

	return &KafkaStorage{
		Producer: producer,
		Consumer: consumer,
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
	}, nil
}

// Produce sends a message to the Kafka topic
func (k *KafkaStorage) Produce(message string) error {
	// Send the message to the topic
	partition, offset, err := k.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: k.Topic,
		Value: sarama.StringEncoder(message),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d", k.Topic, partition, offset)
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
			err := k.Consumer.Consume(ctx, []string{k.Topic}, &consumerGroupHandler{handler: handler})
			if err != nil {
				log.Printf("Error while consuming message: %v", err)
				time.Sleep(1 * time.Second)
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
	// Close the producer
	if err := k.Producer.Close(); err != nil {
		log.Printf("Error while closing Kafka producer: %v", err)
	}

	// Close the consumer
	if err := k.Consumer.Close(); err != nil {
		log.Printf("Error while closing Kafka consumer: %v", err)
	}
}
