package rulesengineservice

import (
	"github.com/case-management-suite/scheduler"
	"go.uber.org/fx"
)

type RulesServiceClientParams struct {
	fx.In
	Scheduler scheduler.WorkScheduler
}

type RulesServiceClientResult struct {
	fx.Out
	Client RulesServiceClient
}

func NewClient(params RulesServiceClientParams) RulesServiceClientResult {
	client := NewRulesServiceClient(params.Scheduler)
	return RulesServiceClientResult{Client: client}
}

var RulesServiceClientModule = fx.Module("rules_engine_client",
	fx.Provide(
		NewClient,
	),
)
