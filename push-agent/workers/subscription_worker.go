package workers

import (
	"github.com/RichardKnop/machinery/v1"
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
		machineryServer *machinery.Server
		subscriptionService services.SubscriptionService
		taskName string
	}
)

func (w *subscriptionWorker) DispatchWorker() error {
	err := w.machineryServer.RegisterTask(w.taskName, w.subscriptionService.HandlePublishTask)
	if err != nil {
		w.logger.Error("failed to register publish task", zap.Error(err))
		return err
	}

	worker := w.machineryServer.NewWorker("publish_worker", 0)
	err = worker.Launch()
	if err != nil {
		w.logger.Error("failed to launch publish worker", zap.Error(err))
		return err
	}

	return nil
}

func NewSubscriptionWorker(config *viper.Viper, logger *zap.Logger, machineryServer *machinery.Server, subscriptionService services.SubscriptionService) SubscriptionWorker {
	enabled := config.GetBool("workers.subscription.enabled")
	taskName := config.GetString("redis.pubsub.publish_task")

	return &subscriptionWorker{
		enabled: enabled,
		logger: logger.Named("subscriptionWorker"),
		machineryServer: machineryServer,
		subscriptionService: subscriptionService,
		taskName: taskName,
	}
}
