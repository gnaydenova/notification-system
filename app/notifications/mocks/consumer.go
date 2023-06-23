package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Consumer struct {
	mock.Mock
}

func (c *Consumer) Consume(ctx context.Context) {
	c.Called(ctx)
}
