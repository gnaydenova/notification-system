package notifications

import (
	"context"
	"encoding/json"

	"github.com/gnaydenova/notification-system/app/notifications/dto"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Produce(ctx context.Context, msg dto.Notification) error
	Close() error
}

type ProducerConfig struct {
	Topic string `yaml:"topic"`
	Addr  []string `yaml:"addr"`
}

type KafkaProducer struct {
	writer *kafka.Writer
	logger kafka.Logger
}

func NewProducer(cfg ProducerConfig, l Logger) Producer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(cfg.Addr...),
			Topic:                  cfg.Topic,
			AllowAutoTopicCreation: true,
			Logger:                 l,
		},
		logger: l,
	}
}

func (p *KafkaProducer) Produce(ctx context.Context, msg dto.Notification) error {
	val, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		// Nil key means the data will be sent in a Round Robin fashion between partitions.
		// Key:   nil,
		Value: val,
	})

	if err != nil {
		p.logger.Printf("error writing message %s\n", err.Error())
	}

	return err
}

func (p *KafkaProducer) Close() error {
	p.logger.Printf("closing writer\n")
	return p.writer.Close()
}
