package rules_test

import (
	"configlinter/internal/domain"
	"configlinter/internal/rules"
	"fmt"
	"testing"
)

func TestPlaintextPassword_Variants(t *testing.T) {
	cases := []struct {
		key     string
		value   any
		trigger bool
	}{
		{"password", "secret123", true},
		{"db_password", "admin", true},
		{"PASSWORD", "test", true},
		{"pass", "test", false},
		{"hostname", "localhost", false},
		{"password", "", false}, // пустой — может не триггерить
	}

	rule := rules.NewPlaintextPasswordRule()
	for _, tt := range cases {
		t.Run(tt.key, func(t *testing.T) {
			node := &domain.ConfigNode{
				Key:   tt.key,
				Path:  "config." + tt.key,
				Value: tt.value,
			}
			root := &domain.ConfigNode{
				Key:        "root",
				Path:       "root",
				Subsidiary: []*domain.ConfigNode{node},
			}

			findings := rule.Analyze(root)
			if tt.trigger && len(findings) == 0 {
				t.Errorf("expected finding for key='%s' value='%v'", tt.key, tt.value)
			}
			if !tt.trigger && len(findings) > 0 {
				t.Errorf("unexpected finding for key='%s' value='%v'", tt.key, tt.value)
			}
		})
	}
}

func TestBindAll_Variants(t *testing.T) {
	cases := []struct {
		key     string
		value   any
		trigger bool
	}{
		{"host", "0.0.0.0", true},
		{"bind", "0.0.0.0", true},
		{"host", "127.0.0.1", false},
		{"host", "localhost", false},
		{"address", "0.0.0.0", true},
		{"port", "0.0.0.0", false},
	}

	rule := rules.NewBindAllRule()
	for _, tt := range cases {
		t.Run(tt.key+"_"+fmt.Sprint(tt.value), func(t *testing.T) {
			node := &domain.ConfigNode{
				Key:   tt.key,
				Path:  "server." + tt.key,
				Value: tt.value,
			}
			root := &domain.ConfigNode{
				Key:        "root",
				Path:       "root",
				Subsidiary: []*domain.ConfigNode{node},
			}

			findings := rule.Analyze(root)
			if tt.trigger && len(findings) == 0 {
				t.Errorf("expected finding for %s=%v", tt.key, tt.value)
			}
			if !tt.trigger && len(findings) > 0 {
				t.Errorf("unexpected finding for %s=%v", tt.key, tt.value)
			}
		})
	}
}

func TestWeakCrypto_Algorithms(t *testing.T) {
	cases := []struct {
		value   string
		trigger bool
	}{
		{"md5", true},
		{"sha1", true},
		{"des", true},
		{"rc4", true},
		{"aes256", false},
		{"sha256", false},
		{"rsa", false},
	}

	rule := rules.NewWeakCryptoRule()
	for _, tt := range cases {
		t.Run(tt.value, func(t *testing.T) {
			node := &domain.ConfigNode{
				Key:   "algorithm",
				Path:  "crypto.algorithm",
				Value: tt.value,
			}
			root := &domain.ConfigNode{
				Key:        "root",
				Path:       "root",
				Subsidiary: []*domain.ConfigNode{node},
			}

			findings := rule.Analyze(root)
			if tt.trigger && len(findings) == 0 {
				t.Errorf("expected finding for '%s'", tt.value)
			}
			if !tt.trigger && len(findings) > 0 {
				t.Errorf("unexpected finding for '%s'", tt.value)
			}
		})
	}
}

func TestTLSDisabled_Values(t *testing.T) {
	cases := []struct {
		key     string
		value   any
		trigger bool
	}{
		{"tls", false, true},
		{"ssl", false, true},
		{"tls_enabled", false, true},
		{"tls", true, false},
		{"tls", "false", true},
		{"debug", false, false},
	}

	rule := rules.NewTLSDisabledRule()
	for _, tt := range cases {
		t.Run(tt.key+"_"+fmt.Sprint(tt.value), func(t *testing.T) {
			node := &domain.ConfigNode{
				Key:   tt.key,
				Path:  "server." + tt.key,
				Value: tt.value,
			}
			root := &domain.ConfigNode{
				Key:        "root",
				Path:       "root",
				Subsidiary: []*domain.ConfigNode{node},
			}

			findings := rule.Analyze(root)
			if tt.trigger && len(findings) == 0 {
				t.Errorf("expected finding for %s=%v", tt.key, tt.value)
			}
			if !tt.trigger && len(findings) > 0 {
				t.Errorf("unexpected finding for %s=%v", tt.key, tt.value)
			}
		})
	}
}

func TestDebugLog_Values(t *testing.T) {
	cases := []struct {
		key     string
		value   any
		trigger bool
	}{
		{"level", "debug", true},
		{"log_level", "DEBUG", true},
		{"level", "info", false},
		{"level", "warn", false},
		{"level", "error", false},
		{"mode", "debug", false},
	}

	rule := rules.NewDebugLogRule()
	for _, tt := range cases {
		t.Run(tt.key+"_"+fmt.Sprint(tt.value), func(t *testing.T) {
			node := &domain.ConfigNode{
				Key:   tt.key,
				Path:  "logging." + tt.key,
				Value: tt.value,
			}
			root := &domain.ConfigNode{
				Key:        "root",
				Path:       "root",
				Subsidiary: []*domain.ConfigNode{node},
			}

			findings := rule.Analyze(root)
			if tt.trigger && len(findings) == 0 {
				t.Errorf("expected finding for %s=%v", tt.key, tt.value)
			}
			if !tt.trigger && len(findings) > 0 {
				t.Errorf("unexpected finding for %s=%v", tt.key, tt.value)
			}
		})
	}
}
