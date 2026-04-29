package rules_test

import (
	"configlinter/internal/domain"
	"configlinter/internal/rules"
	"testing"
)

func TestWeakCryptoRule_Triggers(t *testing.T) {
	algos := []string{"md5", "sha1", "des", "3des", "rc4", "rc2", "MD5", "SHA1", "DES", "RC4"}

	rule := rules.NewWeakCryptoRule()

	for _, algo := range algos {
		t.Run(algo, func(t *testing.T) {
			root := &domain.ConfigNode{
				Key:  "security",
				Path: "security",
				Subsidiary: []*domain.ConfigNode{
					{Key: "algorithm", Path: "security.algorithm", Value: algo},
				},
			}
			root.Subsidiary[0].Parent = root

			findings := rule.Analyze(root)
			if len(findings) != 1 {
				t.Errorf("algo %q: expected 1 finding, got %d", algo, len(findings))
			}
		})
	}
}

func TestWeakCryptoRule_NoTrigger(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"sha256", "sha256"},
		{"aes-256", "aes-256"},
		{"chacha20", "chacha20"},
		{"argon2", "argon2"},
		{"int value", 256},
		{"bool value", true},
	}

	rule := rules.NewWeakCryptoRule()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := &domain.ConfigNode{
				Key:  "security",
				Path: "security",
				Subsidiary: []*domain.ConfigNode{
					{Key: "algorithm", Path: "security.algorithm", Value: tt.value},
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
