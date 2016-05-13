package server

import (
	"github.com/denderello/ping-pong-grpc/pingpong"

	"golang.org/x/net/context"
)

func (s *GRPCServer) SendPing(ctx context.Context, in *pingpong.Ping) (*pingpong.Pong, error) {
	s.metrics.RPCCounter.Add(1)
	return &pingpong.Pong{Message: "pong"}, nil
}
