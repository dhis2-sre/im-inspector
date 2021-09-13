package handler

import (
	v1 "k8s.io/api/core/v1"
	"log"
)

func ProvideIdHandler() PodHandler {
	return idHandler{}
}

type idHandler struct {
}

func (T idHandler) Supports() string {
	return "dhis2-id"
}

func (T idHandler) Handle(pod v1.Pod) error {
	log.Printf("Id handler invoked: %s", pod.Name)
	return nil
}
