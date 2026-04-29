package integration_test

import (
	"bytes"
	"configlinter/internal/engine"
	"configlinter/internal/parser"
	"configlinter/internal/reporter"
	"configlinter/internal/rules"
	"testing"
)

func TestE2E_JSONConfig(t *testing.T) {
	input := []byte(`{
		"database": {
			"password": "admin123",
			"host": "0.0.0.0",
			"tls": false
		},
		"logging": {
			"level": "debug"
		},
		"crypto": {
			"algorithm": "md5"
		}
	}`)

	p := parser.NewJSONParser()
	root, err := p.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	e := engine.New(
		rules.NewPlaintextPasswordRule(),
		rules.NewBindAllRule(),
		rules.NewTLSDisabledRule(),
		rules.NewDebugLogRule(),
		rules.NewWeakCryptoRule(),
	)

	findings := e.Analyze(root)
	if len(findings) == 0 {
		t.Error("expected findings, got 0")
	}

	var buf bytes.Buffer
	r := reporter.NewJSONReporter()
	if err := r.Report(&buf, findings); err != nil {
		t.Fatalf("reporter error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty report")
	}
}

func TestE2E_CleanConfig(t *testing.T) {
	input := []byte(`{
		"server": {
			"host": "localhost",
			"port": 8080,
			"tls": true
		},
		"logging": {
			"level": "info"
		},
		"crypto": {
			"algorithm": "aes256"
		}
	}`)

	p := parser.NewJSONParser()
	root, err := p.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	e := engine.New(
		rules.NewPlaintextPasswordRule(),
		rules.NewBindAllRule(),
		rules.NewTLSDisabledRule(),
		rules.NewDebugLogRule(),
		rules.NewWeakCryptoRule(),
	)

	findings := e.Analyze(root)
	if len(findings) != 0 {
		for _, f := range findings {
			t.Logf("unexpected: %s at %s: %s", f.RuleID, f.Path, f.Message)
		}
		t.Errorf("expected 0 findings for clean config, got %d", len(findings))
	}
}

func TestE2E_AllFormats(t *testing.T) {
	tests := []struct {
		name   string
		parser parser.Parser
		input  []byte
	}{
		{
			name:   "JSON",
			parser: parser.NewJSONParser(),
			input:  []byte(`{"database": {"password": "secret123"}}`),
		},
		{
			name:   "YAML",
			parser: parser.NewYAMLParser(),
			input:  []byte("database:\n  password: secret123\n"),
		},
		{
			name:   "TOML",
			parser: parser.NewTOMLParser(),
			input:  []byte("[database]\npassword = \"secret123\"\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := tt.parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			e := engine.New(rules.NewPlaintextPasswordRule())
			findings := e.Analyze(root)

			if len(findings) == 0 {
				t.Error("expected plaintext password finding")
			}
		})
	}
}

func TestE2E_TextReporter(t *testing.T) {
	input := []byte(`{
		"server": {"host": "0.0.0.0"},
		"tls": false
	}`)

	p := parser.NewJSONParser()
	root, err := p.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	e := engine.New(
		rules.NewBindAllRule(),
		rules.NewTLSDisabledRule(),
	)

	findings := e.Analyze(root)

	var buf bytes.Buffer
	r := reporter.NewTextReporter()
	if err := r.Report(&buf, findings); err != nil {
		t.Fatalf("reporter error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty text report")
	}
}
