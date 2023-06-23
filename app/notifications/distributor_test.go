package notifications

import (
	"context"
	"testing"

	"github.com/gnaydenova/notification-system/app/channels"
	"github.com/gnaydenova/notification-system/app/notifications/dto"
	"github.com/stretchr/testify/mock"
)

type mockChannel struct {
	mock.Mock
}

func (r *mockChannel) Send(msg string) error {
	args := r.Called(msg)
	return args.Error(0)
}

func TestDistributor(t *testing.T) {
	t.Run("can distribute message through channel", func(t *testing.T) {
		cr := channels.NewRegisrty()
		mockChan := &mockChannel{}
		cr.Add("test", mockChan)

		distributor := NewDistributor(DistributorConfig{}, cr, EmptyLogger)

		msg := "test message"
		mockChan.On("Send", msg).Return(nil)

		distributor.Distribute(context.Background(), dto.Notification{
			Channel: "test",
			Message: msg,
		})

		mockChan.AssertCalled(t, "Send", msg)
		mockChan.AssertNumberOfCalls(t, "Send", 1)
	})
}
