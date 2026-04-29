package reporter_test

import (
	"bytes"
	"configlinter/internal/domain"
	"configlinter/internal/reporter"
	"encoding/json"
	"testing"
)

type jsonOutput struct {
	Total    int              `json:"total"`
	Findings []domain.Finding `json:"findings"`
}

func TestJSONReporter_EmptyFindings(t *testing.T) {
	r := reporter.NewJSONReporter()
	var buf bytes.Buffer

	err := r.Report(&buf, []domain.Finding{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if out.Total != 0 {
		t.Errorf("expected total 0, got %d", out.Total)
	}
	if len(out.Findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(out.Findings))
	}
}

func TestJSONReporter_SingleFinding(t *testing.T) {
	r := reporter.NewJSONReporter()
	var buf bytes.Buffer

	findings := []domain.Finding{
		{RuleID: "R001", Path: "db.password", Severity: domain.HIGH, Message: "test", Recomendation: "fix it"},
	}

	err := r.Report(&buf, findings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if out.Total != 1 {
		t.Errorf("expected total 1, got %d", out.Total)
	}
	if out.Findings[0].RuleID != "R001" {
		t.Errorf("expected RuleID R001, got %s", out.Findings[0].RuleID)
	}
}

func TestJSONReporter_SortBySeverity(t *testing.T) {
	r := reporter.NewJSONReporter()
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

	var out jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if out.Findings[0].Severity != domain.HIGH {
		t.Errorf("expected first HIGH, got %d", out.Findings[0].Severity)
	}
	if out.Findings[1].Severity != domain.MEDIUM {
		t.Errorf("expected second MEDIUM, got %d", out.Findings[1].Severity)
	}
	if out.Findings[2].Severity != domain.LOW {
		t.Errorf("expected third LOW, got %d", out.Findings[2].Severity)
	}
}

func TestJSONReporter_ValidJSON(t *testing.T) {
	r := reporter.NewJSONReporter()
	var buf bytes.Buffer

	findings := []domain.Finding{
		{RuleID: "R001", Path: "a.b", Severity: domain.HIGH, Message: "msg", Recomendation: "rec"},
		{RuleID: "R002", Path: "c.d", Severity: domain.LOW, Message: "msg2", Recomendation: "rec2"},
	}

	err := r.Report(&buf, findings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !json.Valid(buf.Bytes()) {
		t.Error("output is not valid JSON")
	}
}

func TestJSONReporter_NilFindings(t *testing.T) {
	r := reporter.NewJSONReporter()
	var buf bytes.Buffer

	err := r.Report(&buf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Total != 0 {
		t.Errorf("expected total 0, got %d", out.Total)
	}
}
