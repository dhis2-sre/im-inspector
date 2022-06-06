//go:build wireinject
// +build wireinject

package di

import (
	"github.com/dhis2-sre/im-inspector/pgk/config"
	"github.com/dhis2-sre/im-inspector/pgk/handler"
	"github.com/dhis2-sre/im-inspector/pgk/inspector"
	"github.com/dhis2-sre/rabbitmq/pgk/queue"
	"github.com/google/wire"
)

type Environment struct {
	Inspector inspector.Inspector
	Handlers  []handler.PodHandler
}

func ProvideEnvironment(inspector inspector.Inspector, handlers []handler.PodHandler) Environment {
	return Environment{inspector, handlers}
}

func GetEnvironment() Environment {
	wire.Build(
		ProvideEnvironment,
		config.ProvideConfiguration,
		ProvideHandlers,
		inspector.ProvideInspector,
	)
	return Environment{}
}

func ProvideHandlers(configuration config.Configuration) []handler.PodHandler {
	producer := queue.ProvideProducer(configuration.RabbitMq.GetUrl())
	return []handler.PodHandler{
		handler.ProvideTTLDestroyHandler(&producer),
		handler.ProvideTTLWarningHandler(&producer),
		handler.ProvideIdHandler(),
	}
}
