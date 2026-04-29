package parser

import (
	"configlinter/internal/domain"
	"fmt"
	"path/filepath"
	"strings"
)

type Parser interface {
	Extensions() []string

	Parse(data []byte) (*domain.ConfigNode, error)
}

type Registry struct {
	parsers map[string]Parser
}

func NewRegistry() *Registry {
	return &Registry{
		parsers: make(map[string]Parser),
	}
}

func (r *Registry) Register(p Parser) {
	for _, ext := range p.Extensions() {
		r.parsers[ext] = p
	}
}

func (r *Registry) GetByFilename(filename string) (Parser, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	p, ok := r.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
	return p, nil
}

func (r *Registry) GetByFormat(format string) (Parser, error) {
	ext := "." + strings.ToLower(strings.TrimPrefix(format, "."))
	p, ok := r.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
	return p, nil
}
