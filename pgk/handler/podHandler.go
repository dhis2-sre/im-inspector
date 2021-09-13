package handler

import (
	v1 "k8s.io/api/core/v1"
)

type PodHandler interface {
	Supports() string
	Handle(pod v1.Pod) error
}
