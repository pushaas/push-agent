.PHONY: build \
	run

########################################
# app
########################################
setup:
	@go get github.com/oxequa/realize

clean:
	@rm -fr ./dist

build: clean
#	@cp ./config/$(ENV).yml ./dist/config.yml
	@go build -o ./dist/push-agent main.go

run:
	@PUSHAGENT_ENV=local go run main.go

watch:
	@PUSHAGENT_ENV=local realize start --run --no-config

########################################
# docker
########################################

# dev
docker-create-network:
	@docker network create -d bridge push-service-network

docker-build-dev:
	@docker build \
		-f Dockerfile-dev \
		-t push-agent:latest \
		.

docker-run-dev:
	@docker run \
		-it \
		-p 9000:9000 \
		push-agent:latest

docker-build-and-run-dev: docker-build-dev docker-run-dev

# prod
docker-build-prod:
	@docker build \
		-f Dockerfile-prod \
		-t rafaeleyng/push-agent:latest \
		.

docker-run-prod:
	@docker run \
		-it \
		-p 9000:9000 \
		rafaeleyng/push-agent:latest

docker-build-and-run-prod: docker-build-prod docker-run-prod

docker-push-prod: docker-build-prod
	@docker push \
		rafaeleyng/push-agent

########################################
# services
########################################
services-up:
	@docker-compose up -d

services-down:
	@docker-compose down
