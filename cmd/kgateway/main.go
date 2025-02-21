package main

import (
	"context"

	"github.com/solo-io/go-utils/log"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/setup"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/probes"
)

func main() {
	ctx := context.Background()

	// Start a server which is responsible for responding to liveness probes
	probes.StartLivenessProbeServer(ctx)

	if err := setup.Main(ctx); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}
