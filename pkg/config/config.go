package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Inspiration: https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f

type Config struct {
	DeployableNamespaces []string
	RabbitMq             rabbitmq
}

func New() Config {
	namespaces := requireEnv("DEPLOYABLE_NAMESPACES")
	return Config{
		DeployableNamespaces: strings.Split(namespaces, ","),
		RabbitMq: rabbitmq{
			Host:     requireEnv("RABBITMQ_HOST"),
			Port:     getEnvAsInt("RABBITMQ_PORT"),
			Username: requireEnv("RABBITMQ_USERNAME"),
			Password: requireEnv("RABBITMQ_PASSWORD"),
		},
	}
}

type rabbitmq struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (r rabbitmq) GetUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.Username, r.Password, r.Host, r.Port)
}

func requireEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Can't find environment varialbe: %s\n", key)
	}
	return value
}

func getEnvAsInt(key string) int {
	valueStr := requireEnv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Can't parse value as integer: %s", err.Error())
	}
	return value
}
