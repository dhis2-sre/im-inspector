package pod

import (
	"log"

	v1 "k8s.io/api/core/v1"
)

type idHandler struct{}

func NewIDHandler() idHandler {
	return idHandler{}
}

func (h idHandler) Supports() string {
	return "im-id"
}

func (h idHandler) Handle(pod v1.Pod) error {
	log.Printf("Id handler invoked: %s", pod.Name)
	return nil
}
