package rulesengineservice

import (
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/metrics"
	"github.com/case-management-suite/queue"
	"github.com/case-management-suite/scheduler"
	"github.com/rs/zerolog"
)

// Client Factory
func NewRulesServiceClient(sh scheduler.WorkScheduler) RulesServiceClient {
	return RulesServiceClientQueue{Scheduler: sh}
}

// Server Factory 1
func NewRulesServer(Rules RulesService,
	Metrics metrics.MetricsService,
	Scheduler scheduler.WorkScheduler, logger zerolog.Logger) RulesServer {
	return RulesServer{Rules: Rules, Metrics: Metrics, Scheduler: Scheduler}
}

// Server Factory 2
func NewRulesServerFromAppConfig(appConfig config.AppConfig) RulesServer {
	rservice := NewRulesService()

	mservice := metrics.NewCaseMetricsService()

	makeQueueService := queue.QueueServiceFactory(config.RabbitMQ)

	qs := makeQueueService(appConfig.RulesServiceConfig.QueueConfig, appConfig.LogConfig)
	if appConfig.RulesServiceConfig.QueueConfig.PurgeOnStart {
		qs.PurgeAllChannels()
	}
	schservice := scheduler.NewWorkScheduler(qs, appConfig)

	return NewRulesServer(rservice, mservice, schservice, appConfig.RulesServiceConfig.Logger)
}
