tag ?= latest
clean-cmd = docker compose down --remove-orphans --volumes

check:
	pre-commit run --all-files --show-diff-on-failure

docker-image:
	IMAGE_TAG=$(tag) docker compose build prod

init:
	direnv allow
	pip install pre-commit
	pre-commit install --install-hooks --overwrite

smoke-test:
	docker compose up -d rabbitmq kubernetes
	sleep 5
	# TODO: Use z.$$ for file name
	docker compose cp kubernetes:/tmp/kubernetes/k3s.yaml ./k3s.yaml
	yq e -i ".clusters[0].cluster.server = \"https://kubernetes:6443\"" ./k3s.yaml
	IMAGE_TAG=$(tag) docker compose up -d prod
	docker compose logs -f prod

build-dev-image:
	IMAGE_TAG=$(tag) docker compose build dev

push-docker-image:
	IMAGE_TAG=$(tag) docker compose push prod

di:
	wire gen ./pkg/di

launch-dev:
	docker compose up dev rabbitmq

build-test:
	docker compose build test

test: clean
	docker compose run --no-deps test
	$(clean-cmd)

clean:
	$(clean-cmd)

.PHONY: build-image check push-image di init build-dev launch-dev build-test test clean
