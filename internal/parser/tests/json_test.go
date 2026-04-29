package parser_test

import (
	"configlinter/internal/parser"
	"testing"
)

func TestJSONParser_Extensions(t *testing.T) {
	p := parser.NewJSONParser()
	exts := p.Extensions()
	if len(exts) != 1 || exts[0] != ".json" {
		t.Errorf("expected [.json], got %v", exts)
	}
}

func TestJSONParser_SimpleObject(t *testing.T) {
	data := []byte(`{"host": "localhost", "port": 8080}`)
	p := parser.NewJSONParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(root.Subsidiary) != 2 {
		t.Errorf("expected 2 children, got %d", len(root.Subsidiary))
	}
}

func TestJSONParser_NestedObject(t *testing.T) {
	data := []byte(`{"server": {"host": "0.0.0.0", "port": 443}}`)
	p := parser.NewJSONParser()

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

func TestJSONParser_Array(t *testing.T) {
	data := []byte(`{"items": [1, 2, 3]}`)
	p := parser.NewJSONParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var items *struct{ sub int }
	for _, c := range root.Subsidiary {
		if c.Key == "items" {
			if len(c.Subsidiary) != 3 {
				t.Errorf("expected 3 array items, got %d", len(c.Subsidiary))
			}
			return
		}
	}
	_ = items
	t.Error("items key not found")
}

func TestJSONParser_ParentSet(t *testing.T) {
	data := []byte(`{"db": {"password": "secret"}}`)
	p := parser.NewJSONParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	db := root.Subsidiary[0]
	if db.Parent != root {
		t.Error("db.Parent should be root")
	}
	if len(db.Subsidiary) > 0 {
		pw := db.Subsidiary[0]
		if pw.Parent != db {
			t.Error("password.Parent should be db")
		}
	}
}

func TestJSONParser_PathBuilding(t *testing.T) {
	data := []byte(`{"server": {"tls": {"enabled": true}}}`)
	p := parser.NewJSONParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	enabled := root.Subsidiary[0].Subsidiary[0].Subsidiary[0]
	if enabled.Path != "server.tls.enabled" {
		t.Errorf("expected path 'server.tls.enabled', got '%s'", enabled.Path)
	}
}

func TestJSONParser_InvalidJSON(t *testing.T) {
	data := []byte(`{invalid json}`)
	p := parser.NewJSONParser()

	_, err := p.Parse(data)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestJSONParser_EmptyObject(t *testing.T) {
	data := []byte(`{}`)
	p := parser.NewJSONParser()

	root, err := p.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(root.Subsidiary) != 0 {
		t.Errorf("expected 0 children, got %d", len(root.Subsidiary))
	}
}
