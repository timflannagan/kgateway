package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/setup"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {

	stats.StartStatsServer()

	port := flag.Int("p", 8082, "port to bind")
	dev := flag.Bool("dev", false, "use memory instead of connecting to real gloo storage")
	flag.Parse()

	debugMode := os.Getenv("DEBUG") == "1"

	apiServerOpts, err := setup.InitOpts()
	if err != nil {
		return err
	}

	ctx := contextutils.WithLogger(context.Background(), "apiserver")
	contextutils.LoggerFrom(ctx).Infof("listening on :%v", *port)

	// Start the api server
	return setup.Setup(ctx, *port, *dev, debugMode, apiServerOpts)
}
