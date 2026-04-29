package engine

import (
	"configlinter/internal/domain"
	"sort"
)

type Engine struct {
	rules []domain.Rule
}

func New(rules ...domain.Rule) *Engine {
	return &Engine{
		rules: rules,
	}
}

func (e *Engine) Analyze(root *domain.ConfigNode) []domain.Finding {
	var findings []domain.Finding

	for _, rule := range e.rules {
		result := rule.Analyze(root)
		findings = append(findings, result...)
	}
	sort.Slice(findings, func(i, j int) bool {
		return findings[i].Severity > findings[j].Severity
	})
	return findings
}
