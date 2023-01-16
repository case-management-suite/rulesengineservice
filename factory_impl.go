package rulesengineservice

import (
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/metrics"
	"github.com/case-management-suite/common/server"
	"github.com/case-management-suite/common/service"
	"github.com/case-management-suite/queue"
	"github.com/case-management-suite/scheduler"
)

// Client Factory
func NewRulesServiceClient(sh scheduler.WorkScheduler) RulesServiceClient {
	return RulesServiceClientQueue{Scheduler: sh}
}

// Server Factory
func NewRulesServerFromAppConfig(appConfig config.AppConfig) server.Server[RulesServer] {
	rservice := NewRulesService()

	mservice := metrics.NewCaseMetricsService()

	makeQueueService := queue.QueueServiceFactory(config.RabbitMQ)

	return server.NewServer(func(su server.ServerUtils) RulesServer {
		serviceUtils := service.NewServiceUtilsFromServerUtils(su)

		qs := makeQueueService(appConfig.RulesServiceConfig.QueueConfig, serviceUtils)
		if appConfig.RulesServiceConfig.QueueConfig.PurgeOnStart {
			qs.PurgeAllChannels()
		}
		schservice := scheduler.NewWorkScheduler(appConfig, qs, serviceUtils)

		// return NewRulesServer(rservice, mservice, schservice, su)
		return RulesServer{Rules: rservice, Metrics: mservice, Scheduler: schservice, ServerUtils: su}
	}, appConfig)
}

func NewRulesEngineServiceServer(appConfig config.AppConfig, rservice RulesService, mservice metrics.MetricsService, schservice scheduler.WorkScheduler, su server.ServerUtils) server.Server[RulesServer] {
	return server.NewServer(func(su server.ServerUtils) RulesServer {
		// return NewRulesServer(rservice, mservice, schservice, su)
		return RulesServer{Rules: rservice, Metrics: mservice, Scheduler: schservice, ServerUtils: su}
	}, appConfig)
}
