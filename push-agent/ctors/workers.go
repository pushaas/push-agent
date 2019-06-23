package ctors

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/services"
	"github.com/rafaeleyng/push-agent/push-agent/workers"
)

func NewSubscriptionWorker(config *viper.Viper, logger *zap.Logger, machineryServer *machinery.Server, subscriptionService services.SubscriptionService) workers.SubscriptionWorker {
	return workers.NewSubscriptionWorker(config, logger, machineryServer, subscriptionService)
}
