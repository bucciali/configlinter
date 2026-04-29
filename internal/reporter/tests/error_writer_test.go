package reporter_test

import (
	"configlinter/internal/domain"
	"configlinter/internal/reporter"
	"errors"
	"testing"
)

type failWriter struct{}

func (f *failWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("disk full")
}

func TestJSONReporter_WriteError(t *testing.T) {
	r := reporter.NewJSONReporter()
	findings := []domain.Finding{
		{RuleID: "R1", Path: "a", Severity: domain.HIGH, Message: "m"},
	}

	err := r.Report(&failWriter{}, findings)
	if err == nil {
		t.Error("expected error on write failure")
	}
}

func TestTextReporter_WriteError(t *testing.T) {
	r := reporter.NewTextReporter()
	findings := []domain.Finding{
		{RuleID: "R1", Path: "a", Severity: domain.HIGH, Message: "m"},
	}

	err := r.Report(&failWriter{}, findings)
	_ = err
}
