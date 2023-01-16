package rulesengineservice

import (
	"context"

	"github.com/case-management-suite/models"
	"github.com/case-management-suite/scheduler"
	"github.com/rs/zerolog/log"

	"github.com/case-management-suite/common/metrics"
	"github.com/case-management-suite/common/server"
)

type RulesServer struct {
	Rules     RulesService
	Metrics   metrics.MetricsService
	Scheduler scheduler.WorkScheduler
	server.ServerUtils
}

func (rs RulesServer) GetName() string {
	return "rules_server"
}

func (rs RulesServer) GetServerConfig() *server.ServerConfig {
	return &server.ServerConfig{
		Type: server.ProcessServerType,
	}
}

func (rs RulesServer) Start(ctx context.Context) error {
	if err := rs.Scheduler.Start(ctx); err != nil {
		return err
	}
	newCtx := context.Background()
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
	}, newCtx)
}

func (rs RulesServer) Stop(_ context.Context) error {
	return rs.Scheduler.Stop()
}

var _ server.Serveable = RulesServer{}
