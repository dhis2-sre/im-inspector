package pod

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/dhis2-sre/rabbitmq/pgk/queue"

	v1 "k8s.io/api/core/v1"
)

const (
	TtlDestroy = "ttl-destroy"
)

func NewTTLDestroyHandler(producer queueProducer) ttlDestroyHandler {
	return ttlDestroyHandler{producer}
}

type queueProducer interface {
	Produce(channel queue.Channel, payload any)
}

type ttlDestroyHandler struct {
	producer queueProducer
}

func (t ttlDestroyHandler) Supports() string {
	return "im-ttl"
}

func (t ttlDestroyHandler) Handle(pod v1.Pod) error {
	log.Printf("TTL handler invoked: %s", pod.Name)

	creationTimestampLabel := pod.Labels["im-creation-timestamp"]
	ttlLabel := pod.Labels["im-ttl"]
	if creationTimestampLabel == "" || ttlLabel == "" {
		return errors.New("no creationTimestamp or TTL found")
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
		id, err := strconv.ParseUint(pod.Labels["im-id"], 10, 64)
		if err != nil {
			return err
		}
		payload := struct{ ID uint }{uint(id)}
		t.producer.Produce(TtlDestroy, payload)
	}

	return nil
}

// ttlBeforeNow return if creation time + ttl < now
// creationTimestampLabel is a unix timestamp in seconds
// ttlLabel is seconds
func (t ttlDestroyHandler) ttlBeforeNow(creationTimestamp, ttl int64) bool {
	ttlTime := time.Unix(creationTimestamp+ttl, 0).UTC()
	return ttlTime.Before(time.Now())
}
