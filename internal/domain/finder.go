package domain

import "fmt"

type Severity int

const (
	LOW Severity = iota
	MEDIUM
	HIGH
)

func (s Severity) ToString() string {
	switch s {
	case LOW:
		return "LOW"
	case MEDIUM:
		return "MEDIUM"
	case HIGH:
		return "HIGH"
	default:
		return "UNKNOWN"
	}

}

type Finding struct {
	RuleID        string   `json:"rule_id"`
	Severity      Severity `json:"severity"`
	Path          string   `json:"path"`
	Message       string   `json:"message"`
	Recomendation string   `json:"recommendation"`
}

func (f Finding) ToString() string {
	return fmt.Sprintf("%-6s [%s] path: %s\n    %s\n    Рекомендация: %s",
		f.Severity.ToString(), f.RuleID, f.Path, f.Message, f.Recomendation)
}
