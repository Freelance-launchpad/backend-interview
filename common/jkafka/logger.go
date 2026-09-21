package jkafka

import (
	"errors"
	"fmt"
	"log/slog"

	kafka "github.com/segmentio/kafka-go"
)

func kafkaErrorLogger(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)

	for _, arg := range args {
		if err, ok := arg.(error); ok && errors.Is(err, kafka.NotLeaderForPartition) {
			slog.Warn(msg)
			return
		}
	}

	slog.Error(msg)
}
