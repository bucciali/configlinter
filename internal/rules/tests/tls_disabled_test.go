package rules_test

import (
	"configlinter/internal/domain"
	"configlinter/internal/rules"
	"testing"
)

func TestTLSDisabledRule_Triggers(t *testing.T) {
	tests := []struct {
		name      string
		parentKey string
		key       string
		value     any
	}{
		{"tls_enabled false", "server", "tls_enabled", false},
		{"ssl_verify false", "server", "ssl_verify", false},
		{"enabled under tls parent", "tls", "enabled", false},
		{"verify under ssl parent", "ssl", "verify", false},
		{"tls_active false", "server", "tls_active", false},
		{"ssl_secure false", "server", "ssl_secure", false},
	}

	rule := rules.NewTLSDisabledRule()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			child := &domain.ConfigNode{
				Key:   tt.key,
				Path:  tt.parentKey + "." + tt.key,
				Value: tt.value,
			}
			root := &domain.ConfigNode{
				Key:        tt.parentKey,
				Path:       tt.parentKey,
				Subsidiary: []*domain.ConfigNode{child},
			}
			child.Parent = root

			findings := rule.Analyze(root)
			if len(findings) != 1 {
				t.Errorf("expected 1 finding, got %d", len(findings))
			}
		})
	}
}

func TestTLSDisabledRule_NoTrigger(t *testing.T) {
	tests := []struct {
		name      string
		parentKey string
		key       string
		value     any
	}{
		{"tls_enabled true", "server", "tls_enabled", true},
		{"ssl_verify true", "server", "ssl_verify", true},
		{"enabled true under tls", "tls", "enabled", true},
		{"unrelated bool", "server", "debug", false},
		{"string value", "tls", "enabled", "false"},
		{"unrelated key", "server", "port", 443},
	}

	rule := rules.NewTLSDisabledRule()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			child := &domain.ConfigNode{
				Key:   tt.key,
				Path:  tt.parentKey + "." + tt.key,
				Value: tt.value,
			}
			root := &domain.ConfigNode{
				Key:        tt.parentKey,
				Path:       tt.parentKey,
				Subsidiary: []*domain.ConfigNode{child},
			}
			child.Parent = root

			findings := rule.Analyze(root)
			if len(findings) != 0 {
				t.Errorf("expected 0 findings, got %d", len(findings))
			}
		})
	}
}
