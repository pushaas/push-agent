package workers

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/services"
)

type (
	StatsWorker interface {
		DispatchWorker()
	}

	statsWorker struct {
		enabled        bool
		expiration     time.Duration
		interval       time.Duration
		logger         *zap.Logger
		agentName      string
		quitChan       chan struct{}
		statsService   services.StatsService
		workersEnabled bool
	}
)

func (w *statsWorker) performAction() {
	go w.statsService.UpdateGlobalStats(w.agentName, w.expiration)
	go w.statsService.UpdateChannelsStats(w.agentName, w.expiration)
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
			w.logger.Info("stopping stats worker")
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

func NewStatsWorker(config *viper.Viper, logger *zap.Logger, agentName string, statsService services.StatsService) StatsWorker {
	enabled := config.GetBool("workers.stats.enabled")
	expiration := config.GetDuration("workers.stats.expiration")
	interval := config.GetDuration("workers.stats.interval")
	workersEnabled := config.GetBool("workers.enabled")

	return &statsWorker{
		enabled:        enabled,
		expiration:     expiration,
		interval:       interval,
		logger:         logger.Named("statsWorker"),
		agentName:      agentName,
		statsService:   statsService,
		workersEnabled: workersEnabled,
	}
}

