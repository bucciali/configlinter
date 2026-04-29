package parser_test

import (
	"configlinter/internal/parser"
	"testing"
)

func TestRegistry_RegisterAndGetByFilename(t *testing.T) {
	r := parser.NewRegistry()
	r.Register(parser.NewJSONParser())
	r.Register(parser.NewYAMLParser())
	r.Register(parser.NewTOMLParser())

	tests := []struct {
		filename string
		wantErr  bool
	}{
		{"config.json", false},
		{"config.yaml", false},
		{"config.yml", false},
		{"config.toml", false},
		{"config.xml", true},
		{"config.ini", true},
		{"config", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			_, err := r.GetByFilename(tt.filename)
			if tt.wantErr && err == nil {
				t.Errorf("expected error for %s", tt.filename)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error for %s: %v", tt.filename, err)
			}
		})
	}
}

func TestRegistry_GetByFormat(t *testing.T) {
	r := parser.NewRegistry()
	r.Register(parser.NewJSONParser())
	r.Register(parser.NewYAMLParser())
	r.Register(parser.NewTOMLParser())

	tests := []struct {
		format  string
		wantErr bool
	}{
		{"json", false},
		{".json", false},
		{"yaml", false},
		{"yml", false},
		{"toml", false},
		{"JSON", false},
		{"xml", true},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			_, err := r.GetByFormat(tt.format)
			if tt.wantErr && err == nil {
				t.Errorf("expected error for %s", tt.format)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error for %s: %v", tt.format, err)
			}
		})
	}
}

func TestRegistry_EmptyRegistry(t *testing.T) {
	r := parser.NewRegistry()

	_, err := r.GetByFilename("config.json")
	if err == nil {
		t.Error("expected error for empty registry")
	}

	_, err = r.GetByFormat("json")
	if err == nil {
		t.Error("expected error for empty registry")
	}
}

func TestRegistry_OverwriteParser(t *testing.T) {
	r := parser.NewRegistry()
	r.Register(parser.NewJSONParser())
	r.Register(parser.NewJSONParser())

	p, err := r.GetByFilename("config.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Error("expected parser, got nil")
	}
}
