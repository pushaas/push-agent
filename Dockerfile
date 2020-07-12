########################################
# stage 1: build
########################################
FROM golang:1.14 as go-builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN rm -fr ./dist && mkdir ./dist
RUN cp ./config/prod.yml ./dist/prod.yml
RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux make build

########################################
# stage 2: run
########################################
FROM alpine:latest

WORKDIR /app

ENV PUSHAGENT_ENV=prod

COPY --from=go-builder /app/dist/push-agent ./push-agent
COPY --from=go-builder /app/config/prod.yml ./config/prod.yml

ENTRYPOINT ["/app/push-agent"]
