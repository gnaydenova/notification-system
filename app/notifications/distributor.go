package notifications

import (
	"context"
	"fmt"

	"github.com/gnaydenova/notification-system/app/channels"
	"github.com/gnaydenova/notification-system/app/notifications/dto"
)

type DistributorConfig struct {
	RetryConfig ProducerConfig `yaml:"retry_queue"`
	DLQConfig   ProducerConfig `yaml:"dlq"`
	MaxRetries  int            `yaml:"max_retries"`
}

type Distributor interface {
	Distribute(ctx context.Context, n dto.Notification)
}

type DefaultDistributor struct {
	channels channels.Regisrty
	logger   Logger
	cfg      DistributorConfig
}

func NewDistributor(cfg DistributorConfig, r channels.Regisrty, l Logger) Distributor {
	return &DefaultDistributor{
		channels: r,
		logger:   l,
		cfg:      cfg,
	}
}

func (d *DefaultDistributor) Distribute(ctx context.Context, n dto.Notification) {
	channel, ok := d.channels[n.Channel]

	var err error
	if ok {
		d.logger.Printf("sending message via channel '%s'", n.Channel)
		err = channel.Send(n.Message)
	} else {
		err = fmt.Errorf("channel '%s' not found", n.Channel)
	}

	// On error send message to retry queue or dead letter queue if max retries reached.
	if err != nil {
		d.logger.Printf("could not send notification: %s\n", err.Error())

		var p Producer
		if n.RetryCount < d.cfg.MaxRetries {
			n.RetryCount++
			// Send to retry queue.
			p = NewProducer(d.cfg.RetryConfig, d.logger)
		} else {
			// Send to dead letter queue.
			p = NewProducer(d.cfg.DLQConfig, d.logger)
		}

		p.Produce(ctx, n)
		p.Close()
	}
}
