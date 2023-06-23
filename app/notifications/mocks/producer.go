package mocks

import (
	"context"

	"github.com/gnaydenova/notification-system/app/notifications/dto"
	"github.com/stretchr/testify/mock"
)

type Producer struct {
	mock.Mock
}

func (p *Producer) Produce(ctx context.Context, msg dto.Notification) error {
	args := p.Called(ctx, msg)
	return args.Error(0)
}

func (p *Producer) Close() error {
	args := p.Called()
	return args.Error(0)
}
