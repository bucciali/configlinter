package parser

import (
	"configlinter/internal/domain"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type YAMLParser struct{}

func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

func (p *YAMLParser) Extensions() []string {
	return []string{".yaml", ".yml"}
}

func (p *YAMLParser) Parse(data []byte) (*domain.ConfigNode, error) {
	var raw any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("yaml parse error: %w", err)
	}

	normalized := normalize(raw)
	root := buildTree("", "", normalized)
	return root, nil
}

func normalize(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, v := range val {
			out[k] = normalize(v)
		}
		return out
	case map[any]any:
		out := make(map[string]any, len(val))
		for k, v := range val {
			out[fmt.Sprintf("%v", k)] = normalize(v)
		}
		return out
	case []any:
		for i, item := range val {
			val[i] = normalize(item)
		}
		return val
	default:
		return val
	}
}

func ParseFile(path string) (*domain.ConfigNode, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл: %w", err)
	}

	p := NewYAMLParser()
	return p.Parse(data)
}
