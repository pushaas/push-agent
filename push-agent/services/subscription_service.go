package services

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/models"
)

type (
	SubscriptionService interface {
		HandlePublishTask(*string)
	}

	subscriptionService struct{
		logger *zap.Logger
		pushStreamService PushStreamService
	}
)

func (s *subscriptionService) HandlePublishTask(payload *string) {
	var message models.Message
	err := json.Unmarshal([]byte(*payload), &message)
	if err != nil {
		s.logger.Error("failed to unmarshal message", zap.String("payload", *payload), zap.Error(err))
	}

	s.pushStreamService.PublishMessage(&message)
}

func NewSubscriptionService(logger *zap.Logger, pushStreamService PushStreamService) SubscriptionService {
	return &subscriptionService{
		logger: logger.Named("subscriptionService"),
		pushStreamService: pushStreamService,
	}
}
