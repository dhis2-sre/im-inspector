package pod

import (
	"context"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type podGetter struct {
	client *kubernetes.Clientset
}

func NewPodGetter() (*podGetter, error) {
	kubeConfigPath, _ := os.LookupEnv("KUBE_CONFIG_FILE")

	restClientConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return nil, err
	}

	return &podGetter{client: client}, nil
}

func (pg *podGetter) Get(namespaces []string) ([]v1.Pod, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: "im=true",
	}

	var pods []v1.Pod
	for _, namespace := range namespaces {
		list, err := pg.client.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
		if err != nil {
			return nil, err
		}
		pods = append(pods, list.Items...)
	}

	return pods, nil
}
