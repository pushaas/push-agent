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
		Subscribe() error
	}

	subscriptionService struct{
		logger *zap.Logger
		machineryServer *machinery.Server
		pushStreamService PushStreamService
		taskName string
	}
)

func (s *subscriptionService) handlePublishTask(payload string) error {
	var message models.Message
	err := json.Unmarshal([]byte(payload), &message)
	if err != nil {
		s.logger.Error("failed to unmarshal message", zap.String("payload", payload), zap.Error(err))
		return err
	}

	s.pushStreamService.PublishMessage(&message)
	return nil
}

func (s *subscriptionService) Subscribe() error {
	err := s.machineryServer.RegisterTask(s.taskName, s.handlePublishTask)
	if err != nil {
		s.logger.Error("failed to register publish task", zap.Error(err))
		return err
	}

	worker := s.machineryServer.NewWorker("publish_worker", 0)
	err = worker.Launch()
	if err != nil {
		s.logger.Error("failed to launch publish worker", zap.Error(err))
		return err
	}

	return nil
}

func NewSubscriptionService(config *viper.Viper, logger *zap.Logger, pushStreamService PushStreamService, machineryServer *machinery.Server) SubscriptionService {
	return &subscriptionService{
		logger: logger.Named("subscriptionService"),
		machineryServer: machineryServer,
		pushStreamService: pushStreamService,
		taskName: config.GetString("redis.pubsub.publish_task"),
	}
}
