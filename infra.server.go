package rulesengineservice

import (
	"context"

	"github.com/case-management-suite/models"
	"github.com/case-management-suite/scheduler"
	"github.com/rs/zerolog/log"

	"github.com/case-management-suite/common/ctxutils"
	"github.com/case-management-suite/common/metrics"
)

type RulesServer struct {
	Rules     RulesService
	Metrics   metrics.MetricsService
	Scheduler scheduler.WorkScheduler
}

func (rs RulesServer) Start(ctx context.Context) error {
	ctx = ctxutils.DecorateContext(ctx, ctxutils.ContextDecoration{Name: "RulesServer"})
	log.Ctx(ctx).Debug().Str("service", "RulesServer").Msg("Starting rules server")

	if err := rs.Scheduler.Start(ctx); err != nil {
		return err
	}
	return rs.Scheduler.ListenForCaseActions(func(caseAction models.CaseAction) error {
		log.Ctx(ctx).Debug().Str("UUID", caseAction.CaseRecordID).Msg("Processing action...")
		rec, err := rs.Rules.ExecuteAction(caseAction.CaseRecord, caseAction.Action)
		if err != nil {
			return err
		}
		err = rs.Scheduler.NotifyCaseUpdate(rec, ctx)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("Failed to notificate case update")
		}
		return nil
	}, ctx)
}

func (rs RulesServer) Stop() {
	rs.Scheduler.Stop()
}
