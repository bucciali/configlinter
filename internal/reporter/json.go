package reporter

import (
	"configlinter/internal/domain"
	"encoding/json"
	"io"
	"sort"
)

type JSONReporter struct{}

func NewJSONReporter() *JSONReporter {
	return &JSONReporter{}
}

type jsonOutput struct {
	Total    int              `json:"total"`
	Findings []domain.Finding `json:"findings"`
}

func (r *JSONReporter) Report(w io.Writer, findings []domain.Finding) error {
	sort.Slice(findings, func(i, j int) bool {
		return severityOrder(findings[i].Severity) < severityOrder(findings[j].Severity)
	})

	out := jsonOutput{
		Total:    len(findings),
		Findings: findings,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
