package engine_test

import (
	"configlinter/internal/domain"
	"configlinter/internal/engine"
	"testing"
)

type mockRule struct {
	findings []domain.Finding
}

func (m *mockRule) Analyze(_ *domain.ConfigNode) []domain.Finding {
	return m.findings
}

func (m *mockRule) ID() string          { return "mock" }
func (m *mockRule) Name() string        { return "mock" }
func (m *mockRule) Description() string { return "mock rule" }

func TestEngine_NoRules(t *testing.T) {
	e := engine.New()
	root := &domain.ConfigNode{Key: "root", Path: "root"}

	findings := e.Analyze(root)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestEngine_SingleRule(t *testing.T) {
	rule := &mockRule{
		findings: []domain.Finding{
			{Path: "a.b", Severity: domain.HIGH, Message: "test"},
		},
	}

	e := engine.New(rule)
	root := &domain.ConfigNode{Key: "root", Path: "root"}

	findings := e.Analyze(root)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Path != "a.b" {
		t.Errorf("expected path a.b, got %s", findings[0].Path)
	}
}

func TestEngine_MultipleRules(t *testing.T) {
	rule1 := &mockRule{
		findings: []domain.Finding{
			{Path: "a", Severity: domain.LOW, Message: "low"},
		},
	}
	rule2 := &mockRule{
		findings: []domain.Finding{
			{Path: "b", Severity: domain.HIGH, Message: "high"},
			{Path: "c", Severity: domain.MEDIUM, Message: "medium"},
		},
	}

	e := engine.New(rule1, rule2)
	root := &domain.ConfigNode{Key: "root", Path: "root"}

	findings := e.Analyze(root)
	if len(findings) != 3 {
		t.Fatalf("expected 3 findings, got %d", len(findings))
	}
}

func TestEngine_SortBySeverityDesc(t *testing.T) {
	rule := &mockRule{
		findings: []domain.Finding{
			{Path: "a", Severity: domain.LOW, Message: "low"},
			{Path: "b", Severity: domain.HIGH, Message: "high"},
			{Path: "c", Severity: domain.MEDIUM, Message: "medium"},
		},
	}

	e := engine.New(rule)
	root := &domain.ConfigNode{Key: "root", Path: "root"}

	findings := e.Analyze(root)
	if len(findings) != 3 {
		t.Fatalf("expected 3 findings, got %d", len(findings))
	}

	if findings[0].Severity != domain.HIGH {
		t.Errorf("expected first finding HIGH, got %d", findings[0].Severity)
	}
	if findings[1].Severity != domain.MEDIUM {
		t.Errorf("expected second finding MEDIUM, got %d", findings[1].Severity)
	}
	if findings[2].Severity != domain.LOW {
		t.Errorf("expected third finding LOW, got %d", findings[2].Severity)
	}
}

func TestEngine_RuleReturnsEmpty(t *testing.T) {
	rule1 := &mockRule{findings: nil}
	rule2 := &mockRule{findings: []domain.Finding{}}
	rule3 := &mockRule{
		findings: []domain.Finding{
			{Path: "x", Severity: domain.MEDIUM, Message: "found"},
		},
	}

	e := engine.New(rule1, rule2, rule3)
	root := &domain.ConfigNode{Key: "root", Path: "root"}

	findings := e.Analyze(root)
	if len(findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(findings))
	}
}
