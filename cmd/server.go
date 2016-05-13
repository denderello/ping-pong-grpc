package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/denderello/ping-pong-grpc/net"
	"github.com/denderello/ping-pong-grpc/server"
	"github.com/denderello/ping-pong-grpc/telemetry"
)

var (
	serverPort          string
	collectorType       string
	collectorStatsDHost string
	collectorStatsDPort string
)

func init() {
	serverCommand.Flags().StringVar(&serverPort, "port", "8080", "Port to listen on connections.")
	serverCommand.Flags().StringVar(&collectorType, "collector-type", "none", "Metrics collector type to use (none, stdout, statsd, log)")
	serverCommand.Flags().StringVar(&collectorStatsDHost, "collector-statsd-host", "localhost", "StatsD metrics collector host")
	serverCommand.Flags().StringVar(&collectorStatsDPort, "collector-statsd-port", "8125", "StatsD metrics collector port")
}

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Run pingpong in server mode",
	Long:  `Run pingpong in server mode and wait for ping message to respond with a pong.`,
	Run: func(cmd *cobra.Command, args []string) {
		l := newLogger()

		mct, err := telemetry.ParseCollectorType(collectorType)
		if err != nil {
			l.Fatal(err)
		}

		collector, err := telemetry.NewCollector(telemetry.CollectorConfiguration{
			Type: mct,
			StatsDAddress: net.NetAddress{
				Host: collectorStatsDHost,
				Port: collectorStatsDPort,
			},
			Logger: l,
		})
		if err != nil {
			l.Fatal(err)
		}

		s := server.NewGRPCServer(server.GRPCServerConfig{
			Logger: l,
			Address: net.NetAddress{
				Port: serverPort,
			},
			Metrics: server.GRPCServerMetrics{
				RPCCounter: collector.NewCounter("ping_pong_grpc.server.rpcs", 1*time.Second),
			},
		})

		sigs := make(chan os.Signal, 1)
		go func() {
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			l.Debugf("Received signal %s", <-sigs)

			s.Stop()
		}()

		err = s.Start()
		if err != nil {
			l.Fatal(err)
		}
	},
}
