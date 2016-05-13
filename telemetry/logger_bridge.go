package telemetry

import (
	"github.com/denderello/ping-pong-grpc/log"
)

type LoggerBridge struct {
	logger log.Logger
}

func (lb *LoggerBridge) Write(p []byte) (n int, err error) {
	lb.logger.Debugf("Telemetry sending metric: %s", string(p[:]))
	return len(p), nil
}
