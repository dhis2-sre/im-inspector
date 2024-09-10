package pod

import (
	"log/slog"
	"strings"

	v1 "k8s.io/api/core/v1"
)

type Inspector struct {
	logger     *slog.Logger
	handlers   []Handler
	namespaces []string
	pods       podGetter
}

type podGetter interface {
	Get(namespaces []string) ([]v1.Pod, error)
}

func NewInspector(logger *slog.Logger, pods podGetter, namespaces []string, handlers ...Handler) Inspector {
	return Inspector{
		logger:     logger,
		pods:       pods,
		handlers:   handlers,
		namespaces: namespaces,
	}
}

func (i Inspector) Inspect() error {
	handlerMap := i.createHandlersByLabelMap()
	i.logger.Info("Handlers loaded", "count", slog.IntValue(len(i.handlers)))

	pods, err := i.pods.Get(i.namespaces)
	if err != nil {
		return err
	}

	i.logger.Info("Inspecting pods", "count", slog.IntValue(len(pods)))
	for _, pod := range pods {
		i.logger.Info("Inspecting pod", "pod", pod.Name)
		for label := range pod.Labels {
			handlers, exists := handlerMap[label]
			if exists && strings.HasPrefix(label, "im-") {
				for _, h := range handlers {
					err := h.Handle(pod)
					if err != nil {
						i.logger.Info(err.Error())
					}
				}
			}
		}
	}
	i.logger.Info("Inspection done!")

	return nil
}

func (i Inspector) createHandlersByLabelMap() map[string][]Handler {
	handlerMap := make(map[string][]Handler)
	for index := 0; index < len(i.handlers); index++ {
		key := i.handlers[index].Supports()
		handlerMap[key] = append(handlerMap[key], i.handlers[index])
	}
	return handlerMap
}
