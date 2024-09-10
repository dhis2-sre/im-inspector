package pod

import (
	"log/slog"

	v1 "k8s.io/api/core/v1"
)

type idHandler struct {
	logger *slog.Logger
}

func NewIDHandler(logger *slog.Logger) idHandler {
	return idHandler{logger}
}

func (h idHandler) Supports() string {
	return "im-id"
}

func (h idHandler) Handle(pod v1.Pod) error {
	h.logger.Info("Id handler invoked on", "pod", pod.Name)
	return nil
}
