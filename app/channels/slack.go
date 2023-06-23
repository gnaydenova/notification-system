package channels

import (
	"github.com/slack-go/slack"
)

type SlackConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Token     string `yaml:"token"`
	ChannelID string `yaml:"channel_id"`
}

type Slack struct {
	config SlackConfig
}

func NewSlack(cfg SlackConfig) Slack {
	return Slack{config: cfg}
}

func (s Slack) Send(msg string) error {
	if !s.config.Enabled {
		return nil
	}

	api := slack.New(s.config.Token)
	_, _, err := api.PostMessage(s.config.ChannelID, slack.MsgOptionText(msg, true))

	return err
}
