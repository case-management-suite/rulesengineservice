package rulesengineservice

import (
	"fmt"
	"reflect"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/factory"
	"github.com/case-management-suite/common/metrics"
	"github.com/case-management-suite/common/server"
	"github.com/case-management-suite/scheduler"
)

type RuleServiceFactory func() RulesService

type RulesEngineServiceClientFactory func(scheduler.WorkScheduler) RulesServiceClient

type RulesEngineServiceServerFactory func(config.AppConfig, RulesService, metrics.MetricsService, scheduler.WorkScheduler, server.ServerUtils) server.Server[RulesServer]

type RulesMicoservice = server.Server[RulesServer]

type RulesServiceFactories struct {
	factory.FactorySet
	WorkSchedulerFactories          scheduler.WorkSchedulerFactories
	MetricsServiceFactory           metrics.MetricsServiceFactory
	RuleServiceFactory              RuleServiceFactory
	RulesEngineServiceServerFactory RulesEngineServiceServerFactory
	ServerUtilsFactory              func() server.ServerUtils
}

func (f RulesServiceFactories) BuildRulesService(appConfig config.AppConfig) (*RulesMicoservice, error) {
	if err := factory.ValidateFactorySet(f); err != nil {
		return nil, fmt.Errorf("factory: %s -> %w;", reflect.TypeOf(f).Name(), err)
	}
	scheduler, err := f.WorkSchedulerFactories.BuildWorkScheduler(appConfig)
	if err != nil {
		return nil, err
	}

	// tt := reflect.TypeOf(f.ServerUtilsFactory)
	// log.Info().Interface("tt", tt).Msg("....")

	r := f.RulesEngineServiceServerFactory(
		appConfig,
		f.RuleServiceFactory(),
		f.MetricsServiceFactory(),
		*scheduler,
		f.ServerUtilsFactory())
	return &r, nil
}
