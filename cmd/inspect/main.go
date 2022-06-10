package main

import (
	"fmt"
	"os"

	"github.com/dhis2-sre/im-inspector/pkg/config"
	"github.com/dhis2-sre/im-inspector/pkg/handler"
	"github.com/dhis2-sre/im-inspector/pkg/inspector"
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

	pg, err := pod.NewPodGetter()
	if err != nil {
		return err
	}
	producer := queue.ProvideProducer(config.RabbitMq.GetUrl())
	inspector := inspector.NewInspector(pg, config.DeployableNamespaces,
		handler.NewTTLDestroyHandler(&producer),
		handler.NewIDHandler(),
	)

	return inspector.Inspect()
}
