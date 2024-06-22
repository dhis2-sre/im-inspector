tag ?= latest
clean-cmd = docker compose down --remove-orphans --volumes

init:
	pip install pre-commit
	pre-commit clean
	pre-commit install --install-hooks --overwrite

	go install github.com/direnv/direnv@latest
	direnv version

	go install golang.org/x/tools/cmd/goimports@latest

	go install golang.org/x/tools/cmd/goimports@latest

check:
	pre-commit run --verbose --all-files --show-diff-on-failure

docker-image:
	IMAGE_TAG=$(tag) docker compose build prod

push-docker-image:
	IMAGE_TAG=$(tag) docker compose push prod

smoke-test:
	docker compose up -d rabbitmq kubernetes
	sleep 5
	IMAGE_TAG=$(tag) docker compose up -d prod

test: clean
	docker compose run --no-deps test
	$(clean-cmd)

clean:
	$(clean-cmd)

.PHONY: init check build-image push-image test clean
