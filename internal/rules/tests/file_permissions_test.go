package rules_test

import (
	"io/fs"
	"testing"

	"configlinter/internal/domain"
	"configlinter/internal/rules"
)

func TestFilePermissionsRule(t *testing.T) {
	rule := &rules.FilePermissionsRule{}

	tests := []struct {
		name     string
		mode     fs.FileMode
		wantSev  domain.Severity
		wantHits int
	}{
		{"safe 600", 0o600, 0, 0},
		{"safe 640", 0o640, 0, 0},
		{"world-readable 644", 0o644, domain.MEDIUM, 1},
		{"world-writable 666", 0o666, domain.HIGH, 1},
		{"world-writable 777", 0o777, domain.HIGH, 1},
		{"no mode set", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := &domain.ConfigNode{
				FilePath: "/tmp/test.json",
				FileMode: tt.mode,
			}
			findings := rule.Analyze(root)
			if len(findings) != tt.wantHits {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantHits)
			}
			if tt.wantHits > 0 && findings[0].Severity != tt.wantSev {
				t.Errorf("got severity %d, want %d", findings[0].Severity, tt.wantSev)
			}
		})
	}
}
