package rules_test

import (
	"configlinter/internal/domain"
	"configlinter/internal/rules"
	"testing"
)

func TestPlaintextPasswordRule_Triggers(t *testing.T) {
	keys := []string{"password", "passwd", "secret", "api_key", "apikey", "token", "private_key", "db_password"}

	rule := rules.NewPlaintextPasswordRule()

	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			root := &domain.ConfigNode{
				Key:  "root",
				Path: "root",
				Subsidiary: []*domain.ConfigNode{
					{Key: key, Path: "root." + key, Value: "s3cret123"},
				},
			}
			root.Subsidiary[0].Parent = root

			findings := rule.Analyze(root)
			if len(findings) != 1 {
				t.Errorf("key %q: expected 1 finding, got %d", key, len(findings))
			}
			if len(findings) == 1 && findings[0].Severity != domain.HIGH {
				t.Errorf("expected HIGH severity")
			}
		})
	}
}

func TestPlaintextPasswordRule_NoTrigger(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value any
	}{
		{"env ref", "password", "${DB_PASSWORD}"},
		{"vault ref", "secret", "vault:secret/data/db"},
		{"ssm ref", "api_key", "ssm:/prod/key"},
		{"arn ref", "token", "arn:aws:secretsmanager:us-east-1:123:secret:prod"},
		{"1password ref", "private_key", "op://vault/item/field"},
		{"empty value", "password", ""},
		{"non-secret key", "username", "admin"},
		{"bool value", "password", true},
		{"int value", "password", 12345},
	}

	rule := rules.NewPlaintextPasswordRule()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := &domain.ConfigNode{
				Key:  "root",
				Path: "root",
				Subsidiary: []*domain.ConfigNode{
					{Key: tt.key, Path: "root." + tt.key, Value: tt.value},
				},
			}
			root.Subsidiary[0].Parent = root

			findings := rule.Analyze(root)
			if len(findings) != 0 {
				t.Errorf("expected 0 findings, got %d", len(findings))
			}
		})
	}
}
