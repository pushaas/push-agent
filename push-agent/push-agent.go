package agent

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/ctors"
)

func runApp(config *viper.Viper, logger *zap.Logger) error {
	//port := config.GetString("server.port")
	logger.Info("oii")

	return nil
}

func Run() {
	app := fx.New(
		fx.Provide(
			ctors.NewViper,
			ctors.NewLogger,

			// services
		),
		fx.Invoke(runApp),
	)

	app.Run()
}
