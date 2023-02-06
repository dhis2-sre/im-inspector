package pod

import (
	"log"
	"strconv"
	"time"

	"github.com/dhis2-sre/rabbitmq/pgk/queue"
	v1 "k8s.io/api/core/v1"
)

const (
	TtlDestroy = "ttl-destroy"
)

type ttlDestroyHandler struct {
	producer *queue.Producer
}

func NewTTLDestroyHandler(producer *queue.Producer) ttlDestroyHandler {
	return ttlDestroyHandler{producer}
}

func (t ttlDestroyHandler) Supports() string {
	return "im-ttl"
}

func (t ttlDestroyHandler) Handle(pod v1.Pod) error {
	log.Printf("TTL handler invoked: %s", pod.Name)

	creationTimestamp := pod.Labels["im-creationTimestamp"]
	ttl := pod.Labels["im-ttl"]
	if creationTimestamp == "" || ttl == "" {
		log.Println("No creationTimestamp or TTL found")
		return nil
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

func (t ttlDestroyHandler) ttlBeforeNow(creationTimestamp string, ttl string) bool {
	creationTimestampInt, err := strconv.ParseInt(creationTimestamp, 10, 64)
	if err != nil {
		log.Println(err)
		return false
	}

	ttlInt, err := strconv.ParseInt(ttl, 10, 64)
	if err != nil {
		log.Println(err)
		return false
	}

	ttlTime := time.Unix(creationTimestampInt+ttlInt, 0).UTC()
	return ttlTime.Before(time.Now())
}
