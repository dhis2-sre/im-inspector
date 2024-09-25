package pod

import (
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"strconv"
	"time"

	"github.com/dhis2-sre/rabbitmq-client/pkg/rabbitmq"

	v1 "k8s.io/api/core/v1"
)

const (
	ttlDestroy = "ttl-destroy"
)

func NewTTLDestroyHandler(logger *slog.Logger, producer queueProducer) ttlDestroyHandler {
	return ttlDestroyHandler{logger, producer}
}

type queueProducer interface {
	Produce(channel rabbitmq.Channel, correlationId string, payload any) error
}

type ttlDestroyHandler struct {
	logger   *slog.Logger
	producer queueProducer
}

func (t ttlDestroyHandler) Supports() string {
	return "im-ttl"
}

func (t ttlDestroyHandler) Handle(pod v1.Pod) error {
	t.logger.Info("TTL handler invoked", "pod", pod.Name)

	creationTimestampLabel := pod.Labels["im-creation-timestamp"]
	if creationTimestampLabel == "" {
		return errors.New("failed to find creationTimestamp label")
	}

	ttlLabel := pod.Labels["im-ttl"]
	if ttlLabel == "" {
		t.logger.Info("no TTL label found")
		return nil
	}

	creationTimestamp, err := strconv.ParseInt(creationTimestampLabel, 10, 64)
	if err != nil {
		return err
	}

	ttl, err := strconv.ParseInt(ttlLabel, 10, 64)
	if err != nil {
		return err
	}

	if t.ttlBeforeNow(creationTimestamp, ttl) {
		id, err := strconv.ParseUint(pod.Labels["im-instance-id"], 10, 64)
		if err != nil {
			return err
		}

		payload := struct{ ID uint }{uint(id)}
		correlationId := uuid.NewString()
		err = t.producer.Produce(ttlDestroy, correlationId, payload)
		if err != nil {
			return err
		}
		t.logger.Info("TTL destroyed", "pod", pod.Name, "namespace", pod.Namespace, "correlationId", correlationId)
	}

	return nil
}

// ttlBeforeNow returns true if the pod has expired according to its time to live.
// creationTimestampLabel is a unix timestamp in seconds.
// ttlLabel is the pods time-to-live in seconds.
func (t ttlDestroyHandler) ttlBeforeNow(creationTimestamp, ttl int64) bool {
	ttlTime := time.Unix(creationTimestamp+ttl, 0).UTC()
	return ttlTime.Before(time.Now())
}
