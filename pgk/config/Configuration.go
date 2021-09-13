package config

import (
	"log"
	"os"
	"strings"
)

// Inspiration: https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f

func ProvideConfiguration() Configuration {
	namespaces := requireEnv("DEPLOYABLE_NAMESPACES")
	return Configuration{
		RabbitMqURL:          requireEnv("RABBITMQ_URL"),
		DeployableNamespaces: strings.Split(namespaces, ","),
	}
}

type Configuration struct {
	RabbitMqURL          string
	DeployableNamespaces []string
}

func requireEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Can't find environment varialbe: %s\n", key)
	}
	return value
}
