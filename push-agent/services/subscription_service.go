package services

import (
	"encoding/json"

	"github.com/RichardKnop/machinery/v1"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/models"
)

type (
	SubscriptionService interface {
		HandlePublishTask(string) error
	}

	subscriptionService struct{
		logger *zap.Logger
		machineryServer *machinery.Server
		pushStreamService PushStreamService
	}
)

func (s *subscriptionService) HandlePublishTask(payload string) error {
	var message models.Message
	err := json.Unmarshal([]byte(payload), &message)
	if err != nil {
		s.logger.Error("failed to unmarshal message", zap.String("payload", payload), zap.Error(err))
		return err
	}

	s.pushStreamService.PublishMessage(&message)
	return nil
}

func NewSubscriptionService(config *viper.Viper, logger *zap.Logger, pushStreamService PushStreamService, machineryServer *machinery.Server) SubscriptionService {
	return &subscriptionService{
		logger: logger.Named("subscriptionService"),
		machineryServer: machineryServer,
		pushStreamService: pushStreamService,
	}
}
