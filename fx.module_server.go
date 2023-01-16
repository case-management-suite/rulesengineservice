package rulesengineservice

import (
	"context"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/metrics"
	"github.com/case-management-suite/common/server"
	"github.com/case-management-suite/queue"
	"github.com/case-management-suite/scheduler"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

type RulesServiceServerParams struct {
	fx.In
	AppConfig config.AppConfig
}

func NewQueueServerFx(lc fx.Lifecycle, params RulesServiceServerParams) (server.Server[RulesServer], error) {
	srv := NewRulesServerFromAppConfig(params.AppConfig)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msg("Starting RulesServer")
			return srv.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return srv.Stop(ctx)
		},
	})
	return srv, nil
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
				fx.Private,
			),
			fx.Provide(NewQueueServerFx),
		),
	)
}

func NewRulesServiceCServer(appConfig config.AppConfig) *fx.App {
	return fx.New(
		FxServerOpts(appConfig),
		fx.Invoke(func(_ server.Server[RulesServer]) {}),
	)
}
