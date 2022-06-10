package pod

import (
	"context"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type client struct {
	k8s *kubernetes.Clientset
}

func NewClient() (*client, error) {
	kubeConfigPath, _ := os.LookupEnv("KUBE_CONFIG_FILE")

	restClientConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	k8s, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return nil, err
	}

	return &client{k8s: k8s}, nil
}

func (c *client) Get(namespaces []string) ([]v1.Pod, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: "im=true",
	}

	var pods []v1.Pod
	for _, namespace := range namespaces {
		list, err := c.k8s.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
		if err != nil {
			return nil, err
		}
		pods = append(pods, list.Items...)
	}

	return pods, nil
}
