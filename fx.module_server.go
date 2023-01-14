package rulesengineservice

import (
	"context"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/metrics"
	"github.com/case-management-suite/queue"
	"github.com/case-management-suite/scheduler"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

type RulesServiceServerParams struct {
	fx.In
	AppConfig config.AppConfig
}

type QueueServerFxOut = interface{}

func NewQueueServerFx(lc fx.Lifecycle, params RulesServiceServerParams) (QueueServerFxOut, error) {
	server := NewRulesServerFromAppConfig(params.AppConfig)
	var channel QueueServerFxOut

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msg("Starting RulesServer")
			err := server.Start(ctx)
			return err
		},
		OnStop: func(_ context.Context) error {
			log.Printf("Stopping app..")
			server.Stop()
			return nil
		},
	})
	return channel, nil
}

type Bah struct{}

func FxServerOpts(appConfig config.AppConfig) fx.Option {
	return fx.Options(
		fx.Module("RulesServiceServer",
			config.FxConfig(appConfig),
			fx.Provide(
				NewRulesService,
				func() Bah {
					return Bah{}
				},
				metrics.NewCaseMetricsService,
				queue.QueueServiceFactory(appConfig.RulesServiceConfig.QueueType),
				scheduler.NewWorkScheduler,
				NewRulesServer,
				fx.Private,
			),
			fx.Provide(NewQueueServerFx),
		),
	)
}

func NewRulesServiceCServer(appConfig config.AppConfig) *fx.App {
	return fx.New(
		FxServerOpts(appConfig),
		fx.Invoke(func(c QueueServerFxOut) {}),
	)
}
