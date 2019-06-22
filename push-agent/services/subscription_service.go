package services

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/models"
)

type (
	SubscriptionService interface {
		Subscribe() error
	}

	subscriptionService struct{
		config *viper.Viper
		logger *zap.Logger
		pushStreamService PushStreamService
		redisClient redis.UniversalClient
	}
)

func (s *subscriptionService) handleMessage(payload string) {
	var message models.Message
	err := json.Unmarshal([]byte(payload), &message)
	if err != nil {
		s.logger.Error("failed to unmarshal message", zap.String("payload", payload), zap.Error(err))
		return
	}

	s.pushStreamService.PublishMessage(&message)
}

func (s *subscriptionService) listenMessagesOn(ch <-chan *redis.Message, handler func(string)) {
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
	messagesCh, err := s.subscribeTo(s.config.GetString("redis.pubsub.messages"))
	if err != nil {
		return err
	}

	go s.listenMessagesOn(messagesCh, s.handleMessage)

	return nil
}

func NewSubscriptionService(config *viper.Viper, logger *zap.Logger, redisClient redis.UniversalClient, pushStreamService PushStreamService) SubscriptionService {
	return &subscriptionService{
		config: config,
		logger: logger.Named("subscriptionService"),
		pushStreamService: pushStreamService,
		redisClient: redisClient,
	}
}
