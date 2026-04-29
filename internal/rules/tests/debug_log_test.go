package rules_test

import (
	"configlinter/internal/domain"
	"configlinter/internal/rules"
	"testing"
)

func TestDebugLogRule_Triggers(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"debug level", "level", "debug"},
		{"trace level", "level", "trace"},
		{"DEBUG uppercase", "level", "DEBUG"},
		{"log_level key", "log_level", "debug"},
		{"logLevel key", "logLevel", "trace"},
	}

	rule := rules.NewDebugLogRule()

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
			if len(findings) != 1 {
				t.Errorf("expected 1 finding, got %d", len(findings))
			}
		})
	}
}

func TestDebugLogRule_NoTrigger(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value any
	}{
		{"info level", "level", "info"},
		{"warn level", "level", "warn"},
		{"error level", "level", "error"},
		{"not a level key", "mode", "debug"},
		{"bool value", "level", true},
		{"empty value", "level", ""},
	}

	rule := rules.NewDebugLogRule()

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
