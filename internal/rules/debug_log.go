package rules

import (
	"configlinter/internal/domain"
	"strings"
)

type DebugLogRule struct{}

func NewDebugLogRule() *DebugLogRule {
	return &DebugLogRule{}
}

func (r *DebugLogRule) ID() string {
	return "debug_log"
}

func (r *DebugLogRule) Description() string {
	return "Обнаружение debug/trace уровня логирования"
}

func (r *DebugLogRule) Analyze(root *domain.ConfigNode) []domain.Finding {
	var findings []domain.Finding

	root.Walk(func(node *domain.ConfigNode) {
		if !isLogLevelKey(node.Key) {
			return
		}

		val, ok := node.StringValue()
		if !ok {
			return
		}

		lower := strings.ToLower(val)
		if lower == "debug" || lower == "trace" {
			findings = append(findings, domain.Finding{
				RuleID:        r.ID(),
				Severity:      domain.LOW,
				Path:          node.Path,
				Message:       "Логирование в " + lower + "-режиме",
				Recomendation: "Поменяйте уровень логирования на info или выше в production",
			})
		}
	})

	return findings
}

func isLogLevelKey(key string) bool {
	lower := strings.ToLower(key)
	return lower == "level" ||
		strings.Contains(lower, "log") && strings.Contains(lower, "level")
}
