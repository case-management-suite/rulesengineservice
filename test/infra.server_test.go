package rulesengineservice_test

import (
	"context"
	"testing"
	"time"

	common "github.com/case-management-suite/common/config"
	"github.com/case-management-suite/models"
	"github.com/case-management-suite/rulesengineservice"
	"github.com/rs/zerolog/log"
)

func TestExecuteAction(t *testing.T) {
	appConfig := common.NewLocalTestAppConfig()

	rserver := rulesengineservice.NewRulesServerFromAppConfig(appConfig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := rserver.Start(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to start rules server")
		t.FailNow()
	}

	defer rserver.Stop()

	tchan := make(chan models.CaseRecord)

	rserver.Scheduler.ListenForCaseUpdates(func(cr models.CaseRecord) error {
		tchan <- cr
		return nil
	}, ctx)

	log.Debug().Msg("Executing action...")

	cr := models.CaseRecord{ID: models.NewCaseRecordUUID(), Status: "NEW_CASE"}
	err := rserver.Scheduler.ExecuteCaseAction(cr, "START", ctx)

	if err != nil {
		log.Error().Err(err).Msg("Failed execute action")
	}

	select {
	case <-tchan:
		return
	case <-time.After(time.Second * 5):
		log.Error().Msg("Timed out before receiving events")
		t.Fail()
	}

}
