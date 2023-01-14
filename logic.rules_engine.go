package rulesengineservice

import (
	"embed"

	"github.com/case-management-suite/models"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

//go:embed rules/rules.grl
var f embed.FS

const RULESET_NAME = "CaseManagementRules"

type RulesEngineService struct {
	KnowledgeLibrary ast.KnowledgeLibrary
	DataContext      ast.IDataContext
}

func NewRulesEngineService() RulesEngineService {
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	drls, _ := f.ReadFile("rules/rules.grl")

	// // Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	bs := pkg.NewBytesResource(drls)
	err := ruleBuilder.BuildRuleFromResource(RULESET_NAME, "0.0.1", bs)
	if err != nil {
		panic(err)
	}

	return RulesEngineService{KnowledgeLibrary: *knowledgeLibrary}
}

func (res *RulesEngineService) LoadFacts(cases []models.CaseContext) {
	dataCtx := ast.NewDataContext()
	for _, caser := range cases {
		err := dataCtx.Add("CaseContext", &caser)
		if err != nil {
			panic(err)
		}
	}
	// vv := dataCtx.Get("CaseContext")
	// log.Printf("vv=%v", vv)

	res.DataContext = dataCtx
}

func (res RulesEngineService) ExecuteRules() {
	knowledgeBase := res.KnowledgeLibrary.NewKnowledgeBaseInstance(RULESET_NAME, "0.0.1")

	eng := engine.NewGruleEngine()

	err := eng.Execute(res.DataContext, knowledgeBase)
	if err != nil {
		panic(err)
	}
}
