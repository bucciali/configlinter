package scanner

import (
	"fmt"
	"strings"
)

func FormatResults(results []FileResult) string {
	var b strings.Builder
	totalFindings := 0

	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(&b, "\n%s\n   Error: %s\n", r.Path, r.Err)
			continue
		}

		if len(r.Findings) == 0 {
			fmt.Fprintf(&b, "\n%s — no issues\n", r.Path)
			continue
		}

		fmt.Fprintf(&b, "\n %s — %d issue(s)\n", r.Path, len(r.Findings))
		for _, f := range r.Findings {
			fmt.Fprintf(&b, "   %-6s [%s] path: %s\n", f.Severity.ToString(), f.RuleID, f.Path)
			fmt.Fprintf(&b, "          %s\n", f.Message)
			fmt.Fprintf(&b, "          Рекомендация: %s\n", f.Recomendation)
		}
		totalFindings += len(r.Findings)
	}

	fmt.Fprintf(&b, "\nTotal: %d file(s) scanned, %d issue(s) found\n", len(results), totalFindings)
	return b.String()
}

func HasFindings(results []FileResult) bool {
	for _, r := range results {
		if len(r.Findings) > 0 {
			return true
		}
	}
	return false
}
