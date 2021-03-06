package pod

import (
	"log"
	"strings"

	v1 "k8s.io/api/core/v1"
)

type Inspector struct {
	handlers   []Handler
	namespaces []string
	pods       podGetter
}

type podGetter interface {
	Get(namespaces []string) ([]v1.Pod, error)
}

func NewInspector(pods podGetter, namespaces []string, handlers ...Handler) Inspector {
	return Inspector{
		pods:       pods,
		handlers:   handlers,
		namespaces: namespaces,
	}
}

func (i Inspector) Inspect() error {
	handlerMap := i.createHandlersByLabelMap()
	log.Printf("Found %d handlers\n", len(i.handlers))

	pods, err := i.pods.Get(i.namespaces)
	if err != nil {
		return err
	}

	log.Printf("Inspecting %d pods\n", len(pods))
	for _, pod := range pods {
		log.Printf("Target: %s", pod.Name)
		for label := range pod.Labels {
			handlers, exists := handlerMap[label]
			if exists && strings.HasPrefix(label, "im-") {
				for _, h := range handlers {
					err := h.Handle(pod)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
	log.Println("Inspection ended!")

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
