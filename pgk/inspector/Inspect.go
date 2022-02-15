package inspector

import (
	"github.com/dhis2-sre/im-inspector/pgk/cluster"
	"github.com/dhis2-sre/im-inspector/pgk/config"
	"github.com/dhis2-sre/im-inspector/pgk/handler"
	"log"
	"strings"
)

type Inspector interface {
	Inspect()
}

func ProvideInspector(handlers []handler.PodHandler, configuration config.Configuration) Inspector {
	return &inspector{handlers, configuration}
}

type inspector struct {
	handlers      []handler.PodHandler
	configuration config.Configuration
}

func (i inspector) Inspect() {
	log.Println("Initializing...")

	handlerMap := i.createHandlersByLabelMap()
	log.Printf("Found %d handlers", len(i.handlers))

	pods, err := cluster.GetPods(i.configuration)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Found %d pods", len(pods))

	log.Println("Inspecting...")
	for _, pod := range pods {
		log.Printf("Target: %s", pod.Name)
		for label := range pod.Labels {
			handlers, exists := handlerMap[label]
			if exists && strings.HasPrefix(label, "dhis2-") {
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
}

func (i inspector) createHandlersByLabelMap() map[string][]handler.PodHandler {
	handlerMap := make(map[string][]handler.PodHandler)
	for index := 0; index < len(i.handlers); index++ {
		key := i.handlers[index].Supports()
		handlerMap[key] = append(handlerMap[key], i.handlers[index])
	}
	return handlerMap
}
