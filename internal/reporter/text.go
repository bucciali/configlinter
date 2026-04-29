package reporter

import (
	"configlinter/internal/domain"
	"fmt"
	"io"
	"sort"
)

type TextReporter struct{}

func NewTextReporter() *TextReporter {
	return &TextReporter{}
}

func (r *TextReporter) Report(w io.Writer, findings []domain.Finding) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintln(w, "✅ Проблем не найдено")
		return err
	}

	// Сортируем: HIGH → MEDIUM → LOW
	sort.Slice(findings, func(i, j int) bool {
		return severityOrder(findings[i].Severity) < severityOrder(findings[j].Severity)
	})

	fmt.Fprintf(w, "Найдено проблем: %d\n\n", len(findings))

	for i, f := range findings {
		icon := severityIcon(f.Severity)
		fmt.Fprintf(w, "%d. %s %s [%s]\n", i+1, icon, f.Severity.ToString(), f.RuleID)
		fmt.Fprintf(w, "   Путь: %s\n", f.Path)
		fmt.Fprintf(w, "   %s\n", f.Message)
		fmt.Fprintf(w, "   → %s\n\n", f.Recomendation)
	}

	return nil
}

func severityOrder(s domain.Severity) int {
	switch s {
	case domain.HIGH:
		return 0
	case domain.MEDIUM:
		return 1
	case domain.LOW:
		return 2
	default:
		return 3
	}
}

func severityIcon(s domain.Severity) string {
	switch s {
	case domain.HIGH:
		return "🔴"
	case domain.MEDIUM:
		return "🟡"
	case domain.LOW:
		return "🔵"
	default:
		return "⚪"
	}
}
