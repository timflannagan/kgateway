package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/setup"
	"github.com/kgateway-dev/kgateway/v2/internal/version"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/probes"
)

func main() {
	var (
		kgatewayVersion bool
		metricsAddr     string
		healthProbeAddr string
		pprofAddr       string
		logLevel        string
	)
	cmd := &cobra.Command{
		Use:   "kgateway",
		Short: "Runs the kgateway controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			if kgatewayVersion {
				fmt.Println(version.String())
				return nil
			}
			probes.StartLivenessProbeServer(cmd.Context())
			s := setup.New(
				setup.WithMetricsAddr(metricsAddr),
				setup.WithHealthProbeAddr(healthProbeAddr),
				setup.WithPprofAddr(pprofAddr),
				setup.WithLogLevel(logLevel),
			)
			return s.Start(cmd.Context())
		},
	}
	cmd.Flags().BoolVarP(&kgatewayVersion, "version", "v", false, "Print the version of kgateway")
	cmd.Flags().StringVarP(&metricsAddr, "metrics-addr", "m", ":9092", "The address to listen on for metrics")
	cmd.Flags().StringVarP(&healthProbeAddr, "health-probe-addr", "p", ":9093", "The address to listen on for health probes")
	cmd.Flags().StringVarP(&pprofAddr, "pprof-addr", "p", "127.0.0.1:9099", "The address to listen on for pprof")
	cmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "The log level to use")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
