package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dhis2-sre/im-inspector/pkg/config"
	"github.com/dhis2-sre/im-inspector/pkg/pod"
	"github.com/dhis2-sre/rabbitmq-client/pkg/rabbitmq"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err) // nolint:errcheck
		os.Exit(1)
	}
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.New()
	if err != nil {
		return err
	}

	pc, err := pod.NewClient()
	if err != nil {
		return err
	}
	producer := rabbitmq.ProvideProducer(logger, cfg.RabbitMq.GetUrl())
	inspector := pod.NewInspector(logger, pc, cfg.DeployableNamespaces,
		pod.NewTTLDestroyHandler(logger, &producer),
	)

	return inspector.Inspect()
}
