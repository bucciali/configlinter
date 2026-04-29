package rules

import (
	"configlinter/internal/domain"
	"strings"
)

type PlaintextPasswordRule struct{}

func NewPlaintextPasswordRule() *PlaintextPasswordRule {
	return &PlaintextPasswordRule{}
}

func (r *PlaintextPasswordRule) ID() string {
	return "plaintext_password"
}

func (r *PlaintextPasswordRule) Description() string {
	return "Обнаружение паролей и секретов в открытом виде"
}

func (r *PlaintextPasswordRule) Analyze(root *domain.ConfigNode) []domain.Finding {
	var findings []domain.Finding

	root.Walk(func(node *domain.ConfigNode) {
		if !isSecretKey(node.Key) {
			return
		}

		val, ok := node.StringValue()
		if !ok || val == "" {
			return
		}

		if isReference(val) {
			return
		}

		findings = append(findings, domain.Finding{
			RuleID:        r.ID(),
			Severity:      domain.HIGH,
			Path:          node.Path,
			Message:       "Секрет в открытом виде в ключе '" + node.Key + "'",
			Recomendation: "Используйте переменные окружения (${VAR}) или менеджер секретов (Vault)",
		})
	})

	return findings
}

func isSecretKey(key string) bool {
	lower := strings.ToLower(key)
	secrets := []string{"password", "passwd", "secret", "api_key", "apikey", "token", "private_key"}
	for _, s := range secrets {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}

func isReference(val string) bool {
	prefixes := []string{"${", "vault:", "ssm:", "arn:", "op://"}
	for _, p := range prefixes {
		if strings.HasPrefix(val, p) {
			return true
		}
	}
	return false
}
