package reporter_test

import (
	"bytes"
	"configlinter/internal/domain"
	"configlinter/internal/reporter"
	"strings"
	"testing"
)

func TestTextReporter_EmptyFindings(t *testing.T) {
	r := reporter.NewTextReporter()
	var buf bytes.Buffer

	err := r.Report(&buf, []domain.Finding{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Проблем не найдено") {
		t.Errorf("expected 'Проблем не найдено', got: %s", output)
	}
}

func TestTextReporter_SingleFinding(t *testing.T) {
	r := reporter.NewTextReporter()
	var buf bytes.Buffer

	findings := []domain.Finding{
		{RuleID: "R001", Path: "db.password", Severity: domain.HIGH, Message: "exposed password", Recomendation: "use env"},
	}

	err := r.Report(&buf, findings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Найдено проблем: 1") {
		t.Errorf("expected 'Найдено проблем: 1', got: %s", output)
	}
	if !strings.Contains(output, "db.password") {
		t.Errorf("expected path 'db.password' in output")
	}
	if !strings.Contains(output, "R001") {
		t.Errorf("expected RuleID 'R001' in output")
	}
	if !strings.Contains(output, "exposed password") {
		t.Errorf("expected message in output")
	}
	if !strings.Contains(output, "use env") {
		t.Errorf("expected recommendation in output")
	}
}

func TestTextReporter_SortBySeverity(t *testing.T) {
	r := reporter.NewTextReporter()
	var buf bytes.Buffer

	findings := []domain.Finding{
		{RuleID: "R003", Path: "a", Severity: domain.LOW, Message: "low"},
		{RuleID: "R001", Path: "b", Severity: domain.HIGH, Message: "high"},
		{RuleID: "R002", Path: "c", Severity: domain.MEDIUM, Message: "medium"},
	}

	err := r.Report(&buf, findings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	highIdx := strings.Index(output, "HIGH")
	medIdx := strings.Index(output, "MEDIUM")
	lowIdx := strings.Index(output, "LOW")

	if highIdx == -1 || medIdx == -1 || lowIdx == -1 {
		t.Fatalf("severity labels not found in output:\n%s", output)
	}
	if highIdx > medIdx {
		t.Error("HIGH should appear before MEDIUM")
	}
	if medIdx > lowIdx {
		t.Error("MEDIUM should appear before LOW")
	}
}

func TestTextReporter_ContainsIcons(t *testing.T) {
	r := reporter.NewTextReporter()
	var buf bytes.Buffer

	findings := []domain.Finding{
		{RuleID: "R1", Path: "a", Severity: domain.HIGH, Message: "m"},
		{RuleID: "R2", Path: "b", Severity: domain.MEDIUM, Message: "m"},
		{RuleID: "R3", Path: "c", Severity: domain.LOW, Message: "m"},
	}

	err := r.Report(&buf, findings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "🔴") {
		t.Error("expected 🔴 icon for HIGH")
	}
	if !strings.Contains(output, "🟡") {
		t.Error("expected 🟡 icon for MEDIUM")
	}
	if !strings.Contains(output, "🔵") {
		t.Error("expected 🔵 icon for LOW")
	}
}

func TestTextReporter_NilFindings(t *testing.T) {
	r := reporter.NewTextReporter()
	var buf bytes.Buffer

	err := r.Report(&buf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Проблем не найдено") {
		t.Errorf("expected 'Проблем не найдено' for nil, got: %s", output)
	}
}

func TestTextReporter_MultipleFindings_Count(t *testing.T) {
	r := reporter.NewTextReporter()
	var buf bytes.Buffer

	findings := []domain.Finding{
		{RuleID: "R1", Path: "a", Severity: domain.HIGH, Message: "m1"},
		{RuleID: "R2", Path: "b", Severity: domain.LOW, Message: "m2"},
		{RuleID: "R3", Path: "c", Severity: domain.MEDIUM, Message: "m3"},
		{RuleID: "R4", Path: "d", Severity: domain.HIGH, Message: "m4"},
	}

	err := r.Report(&buf, findings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Найдено проблем: 4") {
		t.Errorf("expected 'Найдено проблем: 4', got: %s", output)
	}
}
