package handler

import (
	"github.com/dhis2-sre/instance-queue/pgk/queue"
	v1 "k8s.io/api/core/v1"
	"log"
	"strconv"
	"time"
)

const (
	TtlWarning = "ttl-warning"
)

func ProvideTTLWarningHandler(producer *queue.Producer) PodHandler {
	return ttlWarningHandler{producer}
}

type ttlWarningHandler struct {
	producer *queue.Producer
}

func (t ttlWarningHandler) Supports() string {
	return "dhis2-ttl"
}

func (t ttlWarningHandler) Handle(pod v1.Pod) error {
	log.Printf("TTL Advice handler invoked: %s", pod.Name)
	ttl := pod.Labels["dhis2-ttl"]
	if ttl != "" &&
		t.shouldWarn(ttl, 48*time.Hour) ||
		t.shouldWarn(ttl, 24*time.Hour) ||
		t.shouldWarn(ttl, 1*time.Hour) {

		id, err := strconv.ParseUint(pod.Labels["dhis2-id"], 10, 64)
		if err != nil {
			return err
		}

		payload := struct {
			ID    uint
			Owner string
		}{
			uint(id),
			pod.Labels["dhis2-owner"],
		}
		t.producer.Produce(TtlWarning, payload)
	} else {
		log.Println("No TTL found")
	}
	return nil
}

func (t ttlWarningHandler) shouldWarn(ttl string, back time.Duration) bool {
	ttlInt, err := strconv.ParseInt(ttl, 10, 64)
	if err != nil {
		log.Println(err)
	}

	ttlUnix := time.Unix(ttlInt, 0)
	loc, _ := time.LoadLocation("UTC")
	ttlUtc := ttlUnix.In(loc)

	now := time.Now()
	interval := 10 * time.Minute
	lower := now.Add(-interval/2 - back)
	upper := now.Add(interval/2 - back)

	return ttlUtc.After(lower) && ttlUtc.Before(upper)
}
