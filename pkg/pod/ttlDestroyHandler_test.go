package pod

import (
	"github.com/dhis2-sre/rabbitmq/pgk/queue"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"testing"
	"time"
)

type mockQueueProducer struct{ mock.Mock }

func (m *mockQueueProducer) Produce(channel queue.Channel, payload any) {
	m.Called(channel, payload)
}

func Test_ttlDestroyHandler_Handle(t *testing.T) {
	producer := &mockQueueProducer{}
	handler := NewTTLDestroyHandler(producer)
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"im-id":                 "1",
				"im-creation-timestamp": strconv.Itoa(int(time.Now().Add(time.Minute * 5).Unix())),
				"im-ttl":                "300",
			},
		},
	}

	err := handler.Handle(pod)

	require.NoError(t, err)
	producer.AssertExpectations(t)
}

func Test_ttlDestroyHandler_Handle_Destroy(t *testing.T) {
	producer := &mockQueueProducer{}
	var channel queue.Channel = "ttl-destroy"
	producer.On("Produce", channel, struct{ ID uint }{ID: 1})
	handler := NewTTLDestroyHandler(producer)
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"im-id":                 "1",
				"im-creation-timestamp": strconv.Itoa(int(time.Now().Add(time.Minute * -5).Unix())),
				"im-ttl":                "300",
			},
		},
	}

	err := handler.Handle(pod)

	require.NoError(t, err)
	producer.AssertExpectations(t)
}
