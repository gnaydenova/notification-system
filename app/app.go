package app

import (
	"log"
	"os"

	"github.com/gnaydenova/notification-system/app/channels"
	"github.com/gnaydenova/notification-system/app/notifications"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Channels struct {
		Email channels.EmailConfig `yaml:"email"`
		Slack channels.SlackConfig `yaml:"slack"`
		SMS   channels.SMSConfig   `yaml:"sms"`
		Log   struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"log"`
	} `yaml:"channels"`
	Producer    notifications.ProducerConfig    `yaml:"producer"`
	Consumer    notifications.ConsumerConfig    `yaml:"consumer"`
	RetryQueue  notifications.ConsumerConfig    `yaml:"retry_queue"`
	DLQ         notifications.ConsumerConfig    `yaml:"dlq"`
	Distributor notifications.DistributorConfig `yaml:"distributor"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func NewRegistryFromConfig(cfg *Config) channels.Regisrty {
	r := channels.NewRegisrty()
	if cfg.Channels.Email.Enabled {
		r.Add(channels.TypeEmail, channels.NewEmail(cfg.Channels.Email))
	}
	if cfg.Channels.Slack.Enabled {
		r.Add(channels.TypeSlack, channels.NewSlack(cfg.Channels.Slack))
	}
	if cfg.Channels.SMS.Enabled {
		r.Add(channels.TypeSMS, channels.NewSMS(cfg.Channels.SMS))
	}
	if cfg.Channels.Log.Enabled {
		r.Add(channels.TypeLog, channels.NewLog(log.New(os.Stdout, "", 0)))
	}

	return r
}
