package rulesengineservice

import (
	"context"

	"github.com/case-management-suite/models"
	"github.com/case-management-suite/scheduler"
	"github.com/rs/zerolog/log"
)

type RulesServiceClient interface {
	ExecuteAction(record models.CaseRecord, action string, context context.Context) error
}

type RulesServiceClientQueue struct {
	Scheduler scheduler.WorkScheduler
}

func (c RulesServiceClientQueue) ExecuteAction(record models.CaseRecord, action string, ctx context.Context) error {
	log.Debug().Str("UUID", record.ID).Str("action", action).Str("service", "RulesServiceClientQueue").Msg("Scheduling action on record")
	return c.Scheduler.ExecuteCaseAction(record, action, ctx)
}
