package notifications

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gnaydenova/notification-system/app/channels"
	"github.com/gnaydenova/notification-system/app/notifications/dto"
	"github.com/gnaydenova/notification-system/app/notifications/mocks"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
)

type mockKafkaReader struct {
	mock.Mock
}

func (r *mockKafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	args := r.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (r *mockKafkaReader) Close() error {
	args := r.Called()
	return args.Error(0)
}

func TestKafkaConsumer(t *testing.T) {
	testCases := []struct {
		name         string
		message      kafka.Message
		notification dto.Notification
		readErr      error
	}{
		{
			name:         "can read and distribute messages",
			message:      kafka.Message{Value: []byte(`{"channel":"sms","message":"test","retry_count":0}`)},
			notification: dto.Notification{Channel: channels.TypeSMS, Message: "test"},
		},
		{
			name:    "does not try to distribute messages on read error",
			readErr: errors.New("read error"),
		},
	}

	newConsumer := func(r *mockKafkaReader, d *mocks.Distributor) *KafkaConsumer {
		consumer := &KafkaConsumer{}
		consumer.reader = r
		consumer.logger = EmptyLogger
		consumer.distributor = d
		return consumer
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockReader := &mockKafkaReader{}
			mockDistributor := &mocks.Distributor{}
			consumer := newConsumer(mockReader, mockDistributor)

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			time.AfterFunc(5*time.Millisecond, cancel)

			mockReader.On("ReadMessage", ctx).Return(tc.message, tc.readErr)
			mockReader.On("Close").Return(nil)
			mockDistributor.On("Distribute", ctx, tc.notification)

			consumer.Consume(ctx)

			mockReader.AssertCalled(t, "ReadMessage", ctx)
			if tc.readErr == nil {
				mockDistributor.AssertCalled(t, "Distribute", ctx, tc.notification)
			} else {
				mockDistributor.AssertNumberOfCalls(t, "Distribute", 0)
			}
			mockReader.AssertCalled(t, "Close")
		})
	}
}
