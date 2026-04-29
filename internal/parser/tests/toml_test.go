package parser_test

import (
	"configlinter/internal/parser"
	"testing"
)

func TestTOMLParser_Extensions(t *testing.T) {
	p := parser.NewTOMLParser()
	exts := p.Extensions()
	if len(exts) != 1 || exts[0] != ".toml" {
		t.Errorf("expected [.toml], got %v", exts)
	}
}

func TestTOMLParser_SimpleObject(t *testing.T) {
	data := []byte("host = \"localhost\"\nport = 8080\n")
	p := parser.NewTOMLParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(root.Subsidiary) != 2 {
		t.Errorf("expected 2 children, got %d", len(root.Subsidiary))
	}
}

func TestTOMLParser_NestedTable(t *testing.T) {
	data := []byte("[server]\nhost = \"0.0.0.0\"\nport = 443\n")
	p := parser.NewTOMLParser()

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

func TestTOMLParser_DeepNesting(t *testing.T) {
	data := []byte("[server.tls]\nenabled = true\n")
	p := parser.NewTOMLParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	server := root.Subsidiary[0]
	if server.Key != "server" {
		t.Fatalf("expected 'server', got '%s'", server.Key)
	}

	tls := server.Subsidiary[0]
	if tls.Key != "tls" {
		t.Fatalf("expected 'tls', got '%s'", tls.Key)
	}

	enabled := tls.Subsidiary[0]
	if enabled.Path != "server.tls.enabled" {
		t.Errorf("expected 'server.tls.enabled', got '%s'", enabled.Path)
	}
}

func TestTOMLParser_InvalidTOML(t *testing.T) {
	data := []byte("[invalid\nkey = ")
	p := parser.NewTOMLParser()

	_, err := p.Parse(data)
	if err == nil {
		t.Error("expected error for invalid TOML")
	}
}

func TestTOMLParser_BoolValues(t *testing.T) {
	data := []byte("debug = true\nverbose = false\n")
	p := parser.NewTOMLParser()

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
