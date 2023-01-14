package rulesengineservice_test

import (
	"testing"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/rulesengineservice"
	"github.com/case-management-suite/testutil"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func TestOpts(t *testing.T) {
	opts := rulesengineservice.FxServerOpts(config.NewLocalTestAppConfig())
	if err := fx.ValidateApp(opts); err != nil {
		log.Error().Err(err).Msg("Failed validation")
		t.Fail()
	}

}

func TestNewRulesServiceCServer(t *testing.T) {
	server := rulesengineservice.NewRulesServiceCServer(config.NewLocalAppConfig())

	testutil.FxAppTest(server, func() {
		err := server.Err()

		if err != nil {
			log.Err(err)
			t.Fail()
		}
	})
}
