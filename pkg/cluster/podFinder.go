package cluster

import (
	"context"
	"log"
	"os"

	"github.com/dhis2-sre/im-inspector/pkg/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetPods(configuration config.Configuration) ([]v1.Pod, error) {
	client := getClient()

	listOptions := metav1.ListOptions{
		LabelSelector: "im=true",
	}

	var pods []v1.Pod
	for _, namespace := range configuration.DeployableNamespaces {
		list, err := client.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
		if err != nil {
			return nil, err
		}
		pods = append(pods, list.Items...)
	}

	return pods, nil
}

func getClient() *kubernetes.Clientset {
	kubeConfigPath, _ := os.LookupEnv("KUBE_CONFIG_FILE")

	restClientConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Println(err)
	}

	client, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		log.Println(err)
	}

	return client
}
