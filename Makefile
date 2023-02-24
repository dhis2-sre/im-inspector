tag ?= latest
clean-cmd = docker compose down --remove-orphans --volumes

init:
	pip install pre-commit
	pre-commit install --install-hooks --overwrite

	go install github.com/direnv/direnv@latest
	direnv version

	go install golang.org/x/tools/cmd/goimports@latest

	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec --version

check:
	pre-commit run --all-files --show-diff-on-failure

docker-image:
	IMAGE_TAG=$(tag) docker compose build prod

push-docker-image:
	IMAGE_TAG=$(tag) docker compose push prod

smoke-test:
	docker compose up -d rabbitmq kubernetes
	sleep 5
	docker compose cp kubernetes:/tmp/kubernetes/k3s.yaml ./k3s.yaml
	yq e -i ".clusters[0].cluster.server = \"https://kubernetes:6443\"" ./k3s.yaml
	IMAGE_TAG=$(tag) docker compose up -d prod
	docker compose logs -f prod

test: clean
	docker compose run --no-deps test
	$(clean-cmd)

clean:
	$(clean-cmd)

.PHONY: init check build-image push-image test clean
