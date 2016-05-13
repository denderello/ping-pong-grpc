package server

import (
	"fmt"
	"net"

	"github.com/denderello/ping-pong-grpc/log"
	lnet "github.com/denderello/ping-pong-grpc/net"
	"github.com/denderello/ping-pong-grpc/pingpong"

	"github.com/go-kit/kit/metrics"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type GRPCServerMetrics struct {
	RPCCounter metrics.Counter
}

type GRPCServerConfig struct {
	Logger  log.Logger
	Address lnet.Addresser
	Metrics GRPCServerMetrics
}

type GRPCServer struct {
	logger  log.Logger
	address lnet.Addresser
	server  *grpc.Server
	metrics GRPCServerMetrics
}

func NewGRPCServer(c GRPCServerConfig) *GRPCServer {
	s := &GRPCServer{
		logger:  c.Logger,
		address: c.Address,
		server:  grpc.NewServer(),
		metrics: c.Metrics,
	}

	pingpong.RegisterPingPongServer(s.server, s)

	return s
}

func (s GRPCServer) Start() error {
	s.logger.Info("Starting server mode")

	lis, err := net.Listen("tcp", s.address.Address())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to listen on, %s", s.address.Address()))
	}

	s.logger.Infof("Listening on %s", s.address.Address())
	s.server.Serve(lis)

	return nil
}

func (s GRPCServer) Stop() {
	s.logger.Info("Stopping server")
	s.server.Stop()
}
