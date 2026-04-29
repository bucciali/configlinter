package rules

import (
	"configlinter/internal/domain"
	"strings"
)

type TLSDisabledRule struct{}

func NewTLSDisabledRule() *TLSDisabledRule {
	return &TLSDisabledRule{}
}

func (r *TLSDisabledRule) ID() string {
	return "tls_disabled"
}

func (r *TLSDisabledRule) Description() string {
	return "Обнаружение отключённого TLS/SSL"
}

func (r *TLSDisabledRule) Analyze(root *domain.ConfigNode) []domain.Finding {
	var findings []domain.Finding

	root.Walk(func(node *domain.ConfigNode) {
		if !isTLSContext(node) {
			return
		}

		disabled := false
		switch v := node.Value.(type) {
		case bool:
			disabled = !v
		case string:
			lower := strings.ToLower(node.Key)
			if containsAny(lower, "tls", "ssl") {
				disabled = strings.EqualFold(v, "false") || v == "0" || strings.EqualFold(v, "no")
			}
		}

		if !disabled {
			return
		}

		findings = append(findings, domain.Finding{
			RuleID:        r.ID(),
			Severity:      domain.HIGH,
			Path:          node.Path,
			Message:       "TLS/SSL отключён: " + node.Key + " = false",
			Recomendation: "Включите TLS и настройте валидные сертификаты",
		})
	})

	return findings
}

func isTLSContext(node *domain.ConfigNode) bool {
	lower := strings.ToLower(node.Key)

	if containsAny(lower, "tls", "ssl") {
		return true
	}

	if hasParentTLS(node) {
		return true
	}

	return false
}

func hasParentTLS(node *domain.ConfigNode) bool {
	current := node.Parent
	for current != nil {
		lower := strings.ToLower(current.Key)
		if strings.Contains(lower, "tls") || strings.Contains(lower, "ssl") {
			return true
		}
		current = current.Parent
	}
	return false
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
