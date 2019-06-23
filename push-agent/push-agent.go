package agent

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/ctors"
	"github.com/rafaeleyng/push-agent/push-agent/workers"
)

func runApp(logger *zap.Logger, statsWorker workers.StatsWorker, subscriptionWorker workers.SubscriptionWorker) error {
	log := logger.Named("runApp")

	// start stats worker
	statsWorker.DispatchWorker()

	// start subscription worker
	err := subscriptionWorker.DispatchWorker()
	if err != nil {
		log.Error("error on dispatching subscriptionWorker", zap.Error(err))
		return err
	}

	return nil
}

func Run() {
	app := fx.New(
		fx.Provide(
			ctors.NewHttpClient,
			ctors.NewReqHttpClient,
			ctors.NewViper,
			ctors.NewLogger,
			ctors.NewRedisClient,
			ctors.NewMachineryServer,

			// services
			ctors.NewPushStreamService,
			ctors.NewStatsService,
			ctors.NewSubscriptionService,

			// workers
			ctors.NewStatsWorker,
			ctors.NewSubscriptionWorker,
		),
		fx.Invoke(runApp),
	)

	app.Run()
}
