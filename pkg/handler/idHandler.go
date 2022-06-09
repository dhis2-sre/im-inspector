package handler

import (
	"log"

	v1 "k8s.io/api/core/v1"
)

type idHandler struct{}

func NewIDHandler() PodHandler {
	return idHandler{}
}

func (T idHandler) Supports() string {
	return "im-id"
}

func (T idHandler) Handle(pod v1.Pod) error {
	log.Printf("Id handler invoked: %s", pod.Name)
	return nil
}
