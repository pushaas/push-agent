package services

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type (
	SubscriptionService interface {
		Subscribe() error
	}

	subscriptionService struct{
		config *viper.Viper
		logger *zap.Logger
		redisClient redis.UniversalClient
	}
)

func handleChannel(payload string) {
	fmt.Println(payload)
}

func handleMessage(payload string) {
	fmt.Println(payload)
}

func handlePublishOn(ch <-chan *redis.Message, handler func(string)) {
	for msg := range ch {
		handler(msg.Payload)
	}
}

func (s *subscriptionService) subscribeTo(channel string) (<-chan *redis.Message, error) {
	pubsub := s.redisClient.Subscribe(channel)
	_, err := pubsub.Receive()
	if err != nil {
		s.logger.Error( "error while subscribing to channel", zap.String("channel", channel), zap.Error(err))
		return nil, err
	}

	return pubsub.Channel(), nil
}

func (s *subscriptionService) Subscribe() error {
	channelsCh, err := s.subscribeTo(s.config.GetString("redis.pubsub.channels"))
	if err != nil {
		return err
	}

	messagesCh, err := s.subscribeTo(s.config.GetString("redis.pubsub.messages"))
	if err != nil {
		return err
	}

	go handlePublishOn(channelsCh, handleChannel)
	go handlePublishOn(messagesCh, handleMessage)

	return nil
}

func NewSubscriptionService(config *viper.Viper, logger *zap.Logger, redisClient redis.UniversalClient) SubscriptionService {
	return &subscriptionService{
		config: config,
		logger: logger.Named("subscriptionService"),
		redisClient: redisClient,
	}
}
