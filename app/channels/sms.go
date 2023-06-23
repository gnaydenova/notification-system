package channels

import (
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	AccountSID string `yaml:"account_sid"`
	Token      string `yaml:"token"`
	To         string `yaml:"to"`
	From       string `yaml:"from"`
}

type SMS struct {
	config SMSConfig
}

func NewSMS(cfg SMSConfig) SMS {
	return SMS{config: cfg}
}

func (s SMS) Send(msg string) error {
	if !s.config.Enabled {
		return nil
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: s.config.AccountSID,
		Password: s.config.Token,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(s.config.To)
	params.SetFrom(s.config.From)
	params.SetBody(msg)

	_, err := client.Api.CreateMessage(params)

	return err
}
