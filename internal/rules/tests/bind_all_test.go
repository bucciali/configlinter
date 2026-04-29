package rules_test

import (
	"configlinter/internal/domain"
	"configlinter/internal/rules"
	"testing"
)

func TestBindAllRule_Triggers(t *testing.T) {
	rule := rules.NewBindAllRule()

	root := &domain.ConfigNode{
		Key:  "server",
		Path: "server",
		Subsidiary: []*domain.ConfigNode{
			{Key: "host", Path: "server.host", Value: "0.0.0.0"},
		},
	}
	root.Subsidiary[0].Parent = root

	findings := rule.Analyze(root)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != domain.MEDIUM {
		t.Errorf("expected MEDIUM severity, got %s", findings[0].Severity.ToString())
	}
}

func TestBindAllRule_NoTrigger(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"localhost", "127.0.0.1"},
		{"specific ip", "192.168.1.1"},
		{"hostname", "myhost.local"},
		{"bool value", true},
		{"int value", 8080},
	}

	rule := rules.NewBindAllRule()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := &domain.ConfigNode{
				Key:  "server",
				Path: "server",
				Subsidiary: []*domain.ConfigNode{
					{Key: "host", Path: "server.host", Value: tt.value},
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
