package agent

import (
	"go.uber.org/fx"

	"github.com/rafaeleyng/push-agent/push-agent/ctors"
	"github.com/rafaeleyng/push-agent/push-agent/workers"
)

func runApp(subscriptionWorker workers.SubscriptionWorker) error {
	err := subscriptionWorker.DispatchWorker()
	return err
}

func Run() {
	app := fx.New(
		fx.Provide(
			ctors.NewViper,
			ctors.NewLogger,
			ctors.NewRedisClient,
			ctors.NewMachineryServer,

			// services
			ctors.NewPushStreamService,
			ctors.NewSubscriptionService,

			// workers
			ctors.NewSubscriptionWorker,
		),
		fx.Invoke(runApp),
	)

	app.Run()
}
