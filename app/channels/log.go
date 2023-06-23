package channels

import "log"

type Log struct {
	logger *log.Logger
}

func NewLog(logger *log.Logger) *Log {
	return &Log{logger: logger}
}

func (l *Log) Send(msg string) error {
	l.logger.Printf("received message: %s\n", msg)
	return nil
}
