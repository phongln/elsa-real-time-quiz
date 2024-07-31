package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type (
	Writer interface {
		ProduceMessageWithKey(ctx context.Context, topic string, key, message []byte)
		ProduceMessage(ctx context.Context, topic string, message []byte)
	}
	Reader interface {
		ReadMessage(ctx context.Context) (kafka.Message, error)
	}

	KafkaWriter struct {
		writer *kafka.Writer
	}
	KafkaReader struct {
		reader *kafka.Reader
	}
)

func GetKafkaWriter(brokers []string) *KafkaWriter {
	return &KafkaWriter{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func GetKafkaReader(brokers []string, topic string, consumerGroup string) *KafkaReader {
	return &KafkaReader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: consumerGroup,
		}),
	}
}

func (k *KafkaWriter) ProduceMessageWithKey(ctx context.Context, topic string, key, message []byte) {
	err := k.writer.WriteMessages(ctx, kafka.Message{
		Value: message,
		Topic: topic,
		Key:   key,
	})
	if err != nil {
		logrus.Error("Error writing message to Kafka:", err)
	}
}

func (k *KafkaWriter) ProduceMessage(ctx context.Context, topic string, message []byte) {
	err := k.writer.WriteMessages(ctx, kafka.Message{
		Value: message,
		Topic: topic,
	})
	if err != nil {
		logrus.Error("Error writing message to Kafka:", err)
	}
}

func (k *KafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return k.reader.ReadMessage(ctx)
}
