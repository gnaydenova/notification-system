package notifications

import (
	"context"
	"encoding/json"
	"io"

	"github.com/gnaydenova/notification-system/app/notifications/dto"
	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Consume(ctx context.Context)
}

type ConsumerConfig struct {
	Topic   string   `yaml:"topic"`
	Addr    []string `taml:"addr"`
	GroupID string   `yaml:"group_id"`
}

type kafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	io.Closer
}

type KafkaConsumer struct {
	reader      kafkaReader
	logger      kafka.Logger
	distributor Distributor
}

func NewConsumer(cfg ConsumerConfig, d Distributor, l Logger) Consumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Addr,
			Topic:   cfg.Topic,
			GroupID: cfg.GroupID,
			// Logger: l,
			ErrorLogger: l,
		}),
		logger:      l,
		distributor: d,
	}
}

func (c *KafkaConsumer) Consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Printf("closing reader\n")
			if err := c.reader.Close(); err != nil {
				c.logger.Printf("failed to close reader: %s\n", err)
			}
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.Printf("could not read message: %s", err.Error())
				continue
			}

			var n dto.Notification
			if err := json.Unmarshal(msg.Value, &n); err != nil {
				c.logger.Printf("could not unmarshal message: %s", err.Error())
				continue
			}

			c.distributor.Distribute(ctx, n)
		}
	}
}
