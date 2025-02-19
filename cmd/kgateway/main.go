package main

import (
	"context"
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
		leaderElection  bool
	)
	cmd := &cobra.Command{
		Use:   "kgateway",
		Short: "Runs the kgateway controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			if kgatewayVersion {
				fmt.Println(version.String())
				return nil
			}
			ctx := context.Background()
			probes.StartLivenessProbeServer(ctx)

			opts := []setup.Option{}
			if leaderElection {
				opts = append(opts, setup.WithLeaderElection())
			}

			controller := setup.New(opts...)
			if err := controller.Start(ctx); err != nil {
				return fmt.Errorf("err in main: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&kgatewayVersion, "version", "v", false, "Print the version of kgateway")
	cmd.Flags().BoolVarP(&leaderElection, "leader-election", "l", false, "Enable leader election")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
