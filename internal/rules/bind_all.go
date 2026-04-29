package rules

import (
	"configlinter/internal/domain"
	"fmt"
	"strings"
)

type BindAllRule struct{}

func NewBindAllRule() *BindAllRule { return &BindAllRule{} }

func (r *BindAllRule) ID() string { return "bind_all_interfaces" }
func (r *BindAllRule) Description() string {
	return "Обнаружение привязки ко всем интерфейсам"
}

func (r *BindAllRule) Analyze(root *domain.ConfigNode) []domain.Finding {
	var findings []domain.Finding

	root.Walk(func(node *domain.ConfigNode) {
		val, ok := node.StringValue()
		if !ok || val != "0.0.0.0" {
			return
		}

		lower := strings.ToLower(node.Key)
		if !containsAny(lower, "host", "bind", "address", "listen", "ip") {
			return
		}

		findings = append(findings, domain.Finding{
			RuleID:        r.ID(),
			Severity:      domain.MEDIUM,
			Path:          node.Path,
			Message:       fmt.Sprintf("Привязка ко всем интерфейсам: %s", val),
			Recomendation: "Используйте 127.0.0.1 или конкретный IP",
		})
	})

	return findings
}
