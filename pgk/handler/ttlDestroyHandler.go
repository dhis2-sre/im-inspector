package handler

import (
	"github.com/dhis2-sre/instance-queue/pgk/queue"
	v1 "k8s.io/api/core/v1"
	"log"
	"strconv"
	"time"
)

const (
	TtlDestroy = "ttl-destroy"
)

func ProvideTTLDestroyHandler(producer *queue.Producer) PodHandler {
	return ttlDestroyHandler{producer}
}

type ttlDestroyHandler struct {
	producer *queue.Producer
}

func (t ttlDestroyHandler) Supports() string {
	return "dhis2-ttl"
}

func (t ttlDestroyHandler) Handle(pod v1.Pod) error {
	log.Printf("TTL handler invoked: %s", pod.Name)
	ttl := pod.Labels["dhis2-ttl"]
	log.Printf("!!!!!TTL: \"%s\"", ttl)
	if ttl != "" && t.ttlBeforeNow(ttl) {
		id, err := strconv.ParseUint(pod.Labels["dhis2-id"], 10, 64)
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
