package ctors

import (
	"strings"

	"github.com/Pallinder/sillyname-go"
	"github.com/RichardKnop/machinery/v1"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/services"
	"github.com/rafaeleyng/push-agent/push-agent/workers"
)

func NewStatsWorker(config *viper.Viper, logger *zap.Logger, machineryServer *machinery.Server, statsService services.StatsService) workers.StatsWorker {
	name := strings.ReplaceAll(strings.ToLower(sillyname.GenerateStupidName()), " ", "-")
	return workers.NewStatsWorker(config, logger, name, machineryServer, statsService)
}

func NewSubscriptionWorker(config *viper.Viper, logger *zap.Logger, machineryServer *machinery.Server, subscriptionService services.SubscriptionService) workers.SubscriptionWorker {
	return workers.NewSubscriptionWorker(config, logger, machineryServer, subscriptionService)
}
