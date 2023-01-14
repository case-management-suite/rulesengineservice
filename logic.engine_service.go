package rulesengineservice

import (
	"github.com/case-management-suite/models"
)

type RulesService interface {
	ExecuteAction(models.CaseRecord, string) (models.CaseRecord, error)
	IsActionSupported(action string) bool
}

type RulesServiceImpl struct {
	Rules RulesEngineService
}

func (res RulesServiceImpl) ExecuteAction(caseRecord models.CaseRecord, action string) (models.CaseRecord, error) {
	context := models.CaseContext{CaseRecord: &caseRecord, Action: action}
	res.Rules.LoadFacts([]models.CaseContext{context})
	res.Rules.ExecuteRules()
	return caseRecord, nil
}

func (RulesServiceImpl) IsActionSupported(action string) bool {
	_, ok := models.BaseSupportedActions[action]
	return ok
}

func NewRulesService() RulesService {
	rules := NewRulesEngineService()
	return RulesServiceImpl{Rules: rules}
}
