package workers

import (
	"time"

	"github.com/RichardKnop/machinery/v1"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/services"
)

type (
	StatsWorker interface {
		DispatchWorker()
	}

	statsWorker struct {
		workersEnabled bool
		enabled bool
		expiration time.Duration
		interval time.Duration
		logger *zap.Logger
		machineryServer *machinery.Server
		name string
		quitChan chan struct{}
		statsService services.StatsService
	}
)

func (w *statsWorker) performAction() {
	go w.statsService.UpdateGlobalStats(w.name, w.expiration)
	go w.statsService.UpdateChannelsStats(w.name, w.expiration)
}

// thanks https://stackoverflow.com/a/16466581/1717979
func (w *statsWorker) startWorker() {
	w.performAction() // run once right away

	ticker := time.NewTicker(w.interval)
	for {
		select {
		case <- ticker.C:
			w.performAction()
		case <- w.quitChan:
			w.quitChan = nil
			ticker.Stop()
			return
		}
	}
}

func (w *statsWorker) stopWorker() {
	if w.quitChan != nil {
		w.quitChan <- struct{}{}
	}
}

func (w *statsWorker) DispatchWorker() {
	if w.workersEnabled && w.enabled {
		go w.startWorker()
	}
}

func NewStatsWorker(config *viper.Viper, logger *zap.Logger, name string, machineryServer *machinery.Server, statsService services.StatsService) StatsWorker {
	workersEnabled := config.GetBool("workers.enabled")
	enabled := config.GetBool("workers.stats.enabled")
	expiration := config.GetDuration("workers.stats.expiration")
	interval := config.GetDuration("workers.stats.interval")

	return &statsWorker{
		workersEnabled: workersEnabled,
		enabled: enabled,
		expiration: expiration,
		interval: interval,
		logger: logger.Named("statsWorker"),
		machineryServer: machineryServer,
		name: name,
		statsService: statsService,
	}
}

