CONTAINER := push-agent
CONTAINER_DEV := push-agent-dev
IMAGE := rafaeleyng/$(CONTAINER)
IMAGE_DEV := rafaeleyng/$(CONTAINER_DEV)
TAG := latest

########################################
# app
########################################
.PHONY: setup
setup:
	@go get github.com/oxequa/realize

.PHONY: clean
clean:
	@rm -fr ./dist

.PHONY: build
build: clean
	@go build -o ./dist/push-agent main.go

.PHONY: run
run:
	@go run main.go

.PHONY: watch
watch:
	@realize start --run --no-config

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
		-t $(IMAGE_DEV):$(TAG) \
		.

.PHONY: docker-run-dev
docker-run-dev: docker-clean-dev
	@docker run \
		-it \
		--name=$(CONTAINER_DEV) \
		--network=host \
		$(IMAGE_DEV):$(TAG)

.PHONY: docker-build-and-run-dev
docker-build-and-run-dev: docker-build-dev docker-run-dev

# prod
.PHONY: docker-clean-prod
docker-clean-prod:
	@-docker rm -f $(CONTAINER)

.PHONY: docker-build-prod
docker-build-prod:
	@docker build \
		-f Dockerfile-prod \
		-t $(IMAGE):$(TAG) \
		.

.PHONY: docker-run-prod
docker-run-prod: docker-clean-prod
	@docker run \
		-it \
		--name=$(CONTAINER) \
		--network=host \
		$(IMAGE):$(TAG)

.PHONY: docker-build-and-run-prod
docker-build-and-run-prod: docker-build-prod docker-run-prod

.PHONY: docker-push-prod
docker-push-prod: docker-build-prod
	@docker push \
		$(IMAGE):$(TAG)
