package mocks

import (
	"context"

	"github.com/gnaydenova/notification-system/app/notifications/dto"
	"github.com/stretchr/testify/mock"
)

type Distributor struct {
	mock.Mock
}

func (d *Distributor) Distribute(ctx context.Context, n dto.Notification) {
	d.Called(ctx, n)
}
