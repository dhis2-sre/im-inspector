package handler

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

func NewTTLDestroyHandler(producer *queue.Producer) PodHandler {
	return ttlDestroyHandler{producer}
}

func (t ttlDestroyHandler) Supports() string {
	return "im-ttl"
}

func (t ttlDestroyHandler) Handle(pod v1.Pod) error {
	log.Printf("TTL handler invoked: %s", pod.Name)
	ttl := pod.Labels["im-ttl"]
	if ttl != "" && t.ttlBeforeNow(ttl) {
		id, err := strconv.ParseUint(pod.Labels["im-id"], 10, 64)
		if err != nil {
			return err
		}
		payload := struct{ ID uint }{uint(id)}
		t.producer.Produce(TtlDestroy, payload)
	} else {
		log.Println("No TTL found")
	}
	return nil
}

func (t ttlDestroyHandler) ttlBeforeNow(ttl string) bool {
	ttlInt, err := strconv.ParseInt(ttl, 10, 64)
	if err != nil {
		log.Println(err)
	}
	ttlTime := time.Unix(ttlInt, 0)
	loc, _ := time.LoadLocation("UTC")
	utc := ttlTime.In(loc)
	now := time.Now()
	return utc.Before(now)
}
