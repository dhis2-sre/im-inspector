package pod

import (
	"errors"
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
	Produce(channel rabbitmq.Channel, payload any) error
}

type ttlDestroyHandler struct {
	logger   *slog.Logger
	producer queueProducer
}

func (t ttlDestroyHandler) Supports() string {
	return "im-ttl"
}

func (t ttlDestroyHandler) Handle(pod v1.Pod) error {
	t.logger.Info("TTL handler invoked on", "pod", pod.Name)

	creationTimestampLabel := pod.Labels["im-creation-timestamp"]
	if creationTimestampLabel == "" {
		return errors.New("no creationTimestamp label found")
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
		err = t.producer.Produce(ttlDestroy, payload)
		if err != nil {
			return err
		}
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
