package telemetry

import (
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/denderello/ping-pong-grpc/log"
	lnet "github.com/denderello/ping-pong-grpc/net"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/statsd"
	"github.com/pkg/errors"
)

const (
	CollectorTypeNone = iota
	CollectorTypeStdout
	CollectorTypeStatsD
	CollectorTypeLog
)

type CollectorType uint8

func ParseCollectorType(name string) (CollectorType, error) {
	switch name {
	case "none":
		return CollectorTypeNone, nil
	case "stdout":
		return CollectorTypeStdout, nil
	case "statsd":
		return CollectorTypeStatsD, nil
	case "log":
		return CollectorTypeLog, nil
	default:
		return 0, errors.Errorf("collector type %s not supported", name)
	}
}

type CollectorConfiguration struct {
	Type          CollectorType
	StatsDAddress lnet.Addresser
	Logger        log.Logger
}

type Collector struct {
	writer io.Writer
}

func NewCollector(c CollectorConfiguration) (*Collector, error) {
	var w io.Writer
	var err error

	switch c.Type {
	case CollectorTypeNone:
		w = ioutil.Discard
	case CollectorTypeStdout:
		w = os.Stdout
	case CollectorTypeStatsD:
		w, err = net.Dial("udp", c.StatsDAddress.Address())
	case CollectorTypeLog:
		w = &LoggerBridge{logger: c.Logger}
	}

	if err != nil {
		return nil, err
	}

	return &Collector{
		writer: w,
	}, nil
}

func (c *Collector) NewCounter(name string, reportInterval time.Duration) metrics.Counter {
	return statsd.NewCounter(c.writer, name, reportInterval)
}
