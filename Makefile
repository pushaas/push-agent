TAG := latest
CONTAINER := push-agent
IMAGE := pushaas/$(CONTAINER)
IMAGE_TAGGED := $(IMAGE):$(TAG)
NETWORK := push-service-network

CONTAINER_DEV := $(CONTAINER)-dev
IMAGE_DEV := pushaas/$(CONTAINER_DEV)
IMAGE_TAGGED_DEV := $(IMAGE_DEV):$(TAG)

########################################
# app
########################################
.PHONY: clean
clean:
	@rm -fr ./dist

.PHONY: build
build: clean
	@go build -o ./dist/push-agent main.go

.PHONY: run
run:
	@go run main.go

.PHONY: kill
kill:
	@-killall push-agent

########################################
# docker
########################################

# dev
.PHONY: docker-clean-dev
docker-clean-dev:
	@-docker rm -f $(CONTAINER_DEV)

.PHONY: docker-build-dev
docker-build-dev:
	@docker build \
		-f Dockerfile-dev \
		-t $(IMAGE_TAGGED_DEV) \
		.

.PHONY: docker-run-dev
docker-run-dev: docker-clean-dev
	@docker run \
		-it \
		--name=$(CONTAINER_DEV) \
		--network=$(NETWORK) \
		$(IMAGE_TAGGED_DEV)

.PHONY: docker-build-and-run-dev
docker-build-and-run-dev: docker-build-dev docker-run-dev

# prod
.PHONY: docker-clean
docker-clean:
	@-docker rm -f $(CONTAINER)

.PHONY: docker-build
docker-build:
	@docker build \
		-f Dockerfile \
		-t $(IMAGE):$(TAG) \
		.

.PHONY: docker-run
docker-run: docker-clean
	@docker run \
		-e PUSHAGENT_PUSH_STREAM__URL="http://push-stream:9080" \
		-e PUSHAGENT_REDIS__URL="redis://push-redis:6379" \
		-it \
		--name=$(CONTAINER) \
		--network=$(NETWORK) \
		$(IMAGE):$(TAG)

.PHONY: docker-build-and-run
docker-build-and-run: docker-build docker-run

.PHONY: docker-push
docker-push: docker-build
	@docker push \
		$(IMAGE):$(TAG)
