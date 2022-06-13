package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Inspiration: https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f

type Config struct {
	DeployableNamespaces []string
	RabbitMq             rabbitmq
}

func New() (Config, error) {
	namespaces, err := requireEnv("DEPLOYABLE_NAMESPACES")
	if err != nil {
		return Config{}, err
	}

	host, err := requireEnv("RABBITMQ_HOST")
	if err != nil {
		return Config{}, err
	}

	port, err := getEnvAsInt("RABBITMQ_PORT")
	if err != nil {
		return Config{}, err
	}

	usr, err := requireEnv("RABBITMQ_USERNAME")
	if err != nil {
		return Config{}, err
	}

	pw, err := requireEnv("RABBITMQ_PASSWORD")
	if err != nil {
		return Config{}, err
	}

	return Config{
		DeployableNamespaces: strings.Split(namespaces, ","),
		RabbitMq: rabbitmq{
			Host:     host,
			Port:     port,
			Username: usr,
			Password: pw,
		},
	}, nil
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

func requireEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("can't find environment varialbe: %s", key)
	}
	return value, nil
}

func getEnvAsInt(key string) (int, error) {
	valueStr, err := requireEnv(key)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("can't parse value as integer: %v", err)
	}
	return value, nil
}
