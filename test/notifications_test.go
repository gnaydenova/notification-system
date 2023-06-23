package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gnaydenova/notification-system/app"
	"github.com/gnaydenova/notification-system/app/channels"
	"github.com/gnaydenova/notification-system/app/notifications"
	"github.com/gnaydenova/notification-system/app/notifications/dto"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
)

const (
	channelTest  = "test"
	channelError = "error"
)

type testChannel struct {
	messages chan string
}

func (c testChannel) Send(msg string) error {
	c.messages <- msg
	return nil
}

type errorChannel struct{}

func (c errorChannel) Send(msg string) error {
	return errors.New("I always return this error")
}

type NotificationsTestSuite struct {
	suite.Suite
	cfg              *app.Config
	teardown         func()
	channel          testChannel
	errChannel       errorChannel
	producer         notifications.Producer
	retryQueueReader *kafka.Reader
}

func TestNotificationsTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationsTestSuite))
}

func (t *NotificationsTestSuite) SetupSuite() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.teardown = cancel

	cfg, err := app.NewConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	t.cfg = cfg

	l := notifications.EmptyLogger
	t.producer = notifications.NewProducer(cfg.Producer, l)

	r := channels.NewRegisrty()
	t.channel.messages = make(chan string)
	r.Add(channelTest, t.channel)
	r.Add(channelError, t.errChannel)

	distributor := notifications.NewDistributor(cfg.Distributor, r, l)
	consumer := notifications.NewConsumer(cfg.Consumer, distributor, l)

	go consumer.Consume(ctx)
}

func (t *NotificationsTestSuite) TearDownSuite() {
	t.producer.Close()
	t.teardown()
}

func (t *NotificationsTestSuite) TestCanDistributeMessage() {
	msg := "test message"
	t.producer.Produce(context.Background(), dto.Notification{Channel: channelTest, Message: msg})

	select {
	case actual := <-t.channel.messages:
		t.Equal(msg, actual)
	case <-time.After(10 * time.Second):
		t.Fail("timeout")
	}
}

func (t *NotificationsTestSuite) TestSendsMessageToRetryQueueWhenErrorOnSendOccurs() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: t.cfg.RetryQueue.Addr,
		Topic:   t.cfg.RetryQueue.Topic,
		GroupID: t.cfg.RetryQueue.GroupID,
	})

	msg := "test message with error"
	t.producer.Produce(context.Background(), dto.Notification{Channel: channelError, Message: msg})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	time.AfterFunc(10*time.Second, cancel)

	actual, err := reader.ReadMessage(ctx)
	reader.Close()
	if err != nil {
		t.Fail(err.Error())
	}

	t.Equal(
		fmt.Sprintf(`{"channel":"%s","message":"%s","retry_count":1}`, channelError, msg),
		string(actual.Value),
	)
}

func (t *NotificationsTestSuite) TestSendsMessageToDLQWhenMaxRetryCountIsExceeded() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: t.cfg.DLQ.Addr,
		Topic:   t.cfg.DLQ.Topic,
		GroupID: t.cfg.DLQ.GroupID,
	})

	maxRetries := t.cfg.Distributor.MaxRetries

	t.producer.Produce(context.Background(), dto.Notification{
		Channel:    "some",
		Message:    "test",
		RetryCount: maxRetries,
	})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	time.AfterFunc(10*time.Second, cancel)

	msg, err := reader.ReadMessage(ctx)
	reader.Close()
	if err != nil {
		t.Fail(err.Error())
	}

	t.Equal(fmt.Sprintf(`{"channel":"some","message":"test","retry_count":%d}`, maxRetries), string(msg.Value))
}
