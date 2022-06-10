package main

import (
	"fmt"
	"os"

	"github.com/dhis2-sre/im-inspector/pkg/config"
	"github.com/dhis2-sre/im-inspector/pkg/pod"
	"github.com/dhis2-sre/rabbitmq/pgk/queue"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := config.New()
	if err != nil {
		return err
	}

	pc, err := pod.NewClient()
	if err != nil {
		return err
	}
	producer := queue.ProvideProducer(config.RabbitMq.GetUrl())
	inspector := pod.NewInspector(pc, config.DeployableNamespaces,
		pod.NewTTLDestroyHandler(&producer),
		pod.NewIDHandler(),
	)

	return inspector.Inspect()
}
