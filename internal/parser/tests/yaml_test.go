package parser_test

import (
	"configlinter/internal/parser"
	"testing"
)

func TestYAMLParser_Extensions(t *testing.T) {
	p := parser.NewYAMLParser()
	exts := p.Extensions()
	if len(exts) != 2 {
		t.Errorf("expected 2 extensions, got %v", exts)
	}
}

func TestYAMLParser_SimpleObject(t *testing.T) {
	data := []byte("host: localhost\nport: 8080\n")
	p := parser.NewYAMLParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(root.Subsidiary) != 2 {
		t.Errorf("expected 2 children, got %d", len(root.Subsidiary))
	}
}

func TestYAMLParser_NestedObject(t *testing.T) {
	data := []byte("server:\n  host: 0.0.0.0\n  port: 443\n")
	p := parser.NewYAMLParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(root.Subsidiary) != 1 {
		t.Fatalf("expected 1 child, got %d", len(root.Subsidiary))
	}

	server := root.Subsidiary[0]
	if server.Key != "server" {
		t.Errorf("expected key 'server', got %s", server.Key)
	}
	if len(server.Subsidiary) != 2 {
		t.Errorf("expected 2 children under server, got %d", len(server.Subsidiary))
	}
}

func TestYAMLParser_Array(t *testing.T) {
	data := []byte("items:\n  - one\n  - two\n  - three\n")
	p := parser.NewYAMLParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, c := range root.Subsidiary {
		if c.Key == "items" {
			if len(c.Subsidiary) != 3 {
				t.Errorf("expected 3 array items, got %d", len(c.Subsidiary))
			}
			return
		}
	}
	t.Error("items key not found")
}

func TestYAMLParser_PathBuilding(t *testing.T) {
	data := []byte("server:\n  tls:\n    enabled: true\n")
	p := parser.NewYAMLParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	enabled := root.Subsidiary[0].Subsidiary[0].Subsidiary[0]
	if enabled.Path != "server.tls.enabled" {
		t.Errorf("expected 'server.tls.enabled', got '%s'", enabled.Path)
	}
}

func TestYAMLParser_InvalidYAML(t *testing.T) {
	data := []byte(":\n  :\n - ][")
	p := parser.NewYAMLParser()

	_, err := p.Parse(data)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestYAMLParser_BoolValues(t *testing.T) {
	data := []byte("debug: true\nverbose: false\n")
	p := parser.NewYAMLParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, c := range root.Subsidiary {
		if c.Key == "debug" {
			if c.Value != true {
				t.Errorf("expected true, got %v", c.Value)
			}
		}
	}
}
