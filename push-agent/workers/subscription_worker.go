package workers

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/services"
)

type (
	SubscriptionWorker interface {
		DispatchWorker() error
	}

	subscriptionWorker struct {
		enabled bool
		logger *zap.Logger
		pubsubChannel string
		redisClient redis.UniversalClient
		subscriptionService services.SubscriptionService
		workersEnabled bool
	}
)

func (w *subscriptionWorker) subscribeTo(channel string) (<-chan *redis.Message, error) {
	pubsub := w.redisClient.Subscribe(channel)
	_, err := pubsub.Receive()

	if err != nil {
		w.logger.Error( "error while subscribing to channel", zap.String("channel", channel), zap.Error(err))
		return nil, err
	}

	return pubsub.Channel(), nil
}

func (w *subscriptionWorker) listenMessagesOn(ch <-chan *redis.Message, handler func(*string)) {
	for msg := range ch {
		go handler(&msg.Payload)
	}
}

func (w *subscriptionWorker) startWorker() error {
	messagesCh, err := w.subscribeTo(w.pubsubChannel)
	if err != nil {
		return err
	}

	go w.listenMessagesOn(messagesCh, w.subscriptionService.HandlePublishTask)

	return nil
}

func (w *subscriptionWorker) DispatchWorker() error {
	if w.workersEnabled && w.enabled {
		 return w.startWorker()
	}
	return nil
}

func NewSubscriptionWorker(config *viper.Viper, logger *zap.Logger, redisClient redis.UniversalClient, subscriptionService services.SubscriptionService) SubscriptionWorker {
	enabled := config.GetBool("workers.subscription.enabled")
	pubsubChannel := config.GetString("redis.pubsub-channels.publish")
	workersEnabled := config.GetBool("workers.enabled")

	return &subscriptionWorker{
		enabled: enabled,
		logger: logger.Named("subscriptionWorker"),
		pubsubChannel: pubsubChannel,
		redisClient: redisClient,
		subscriptionService: subscriptionService,
		workersEnabled: workersEnabled,
	}
}
