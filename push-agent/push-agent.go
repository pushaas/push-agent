package agent

import (
	"go.uber.org/fx"

	"github.com/rafaeleyng/push-agent/push-agent/ctors"
	"github.com/rafaeleyng/push-agent/push-agent/services"
)

func runApp(subscriptionService services.SubscriptionService) error {
	err := subscriptionService.Subscribe()
	return err
}

func Run() {
	app := fx.New(
		fx.Provide(
			ctors.NewViper,
			ctors.NewLogger,
			ctors.NewRedis,

			// services
			ctors.NewPushStreamService,
			ctors.NewSubscriptionService,
		),
		fx.Invoke(runApp),
	)

	app.Run()
}
