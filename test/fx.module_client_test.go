package rulesengineservice_test

import (
	"testing"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/queue"
	"github.com/case-management-suite/rulesengineservice"
	"github.com/case-management-suite/scheduler"
	"github.com/case-management-suite/testutil"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func TestRulesServiceClientModule(t *testing.T) {
	appConfig := config.NewLocalTestAppConfig()
	server := fx.New(
		config.FxConfig(appConfig),
		fx.Supply(appConfig.RulesServiceConfig.QueueConfig),
		fx.Provide(queue.QueueServiceFactory(appConfig.RulesServiceConfig.QueueType)),
		fx.Provide(scheduler.NewWorkScheduler),
		rulesengineservice.RulesServiceClientModule,
		fx.Invoke(func(client rulesengineservice.RulesServiceClient) {}))

	testutil.FxAppTest(server, func() {
		err := server.Err()

		if err != nil {
			log.Err(err)
			t.Fail()
		}
	})
}
