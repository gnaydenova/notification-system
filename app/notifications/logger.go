package notifications

import "github.com/segmentio/kafka-go"

type Logger interface {
	Printf(string, ...interface{})
}

var EmptyLogger = kafka.LoggerFunc(func(s string, i ...interface{}) {})
