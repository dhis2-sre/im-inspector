tag ?= latest
clean-cmd = docker compose down --remove-orphans --volumes

prod-image:
	IMAGE_TAG=$(tag) docker compose build prod

build-dev-image:
	IMAGE_TAG=$(tag) docker compose build dev

push-image:
	IMAGE_TAG=$(tag) docker compose push prod

di:
	wire gen ./pgk/di

launch-dev:
	docker compose up dev rabbitmq

build-test:
	docker compose build test

test: clean
	docker compose run --no-deps test
	$(clean-cmd)

clean:
	$(clean-cmd)

keys:
	echo "no keys needed"

.PHONY: build-image push-image di build-dev launch-dev build-test test clean keys
