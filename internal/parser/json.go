package parser

import (
	"configlinter/internal/domain"
	"encoding/json"
	"fmt"
)

type JSONParser struct{}

func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

func (p *JSONParser) Extensions() []string {
	return []string{".json"}
}

func (p *JSONParser) Parse(data []byte) (*domain.ConfigNode, error) {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("json parse error: %w", err)
	}
	root := buildTree("", "", raw)
	return root, nil
}

func buildTree(key, path string, value any) *domain.ConfigNode {
	node := &domain.ConfigNode{
		Key:  key,
		Path: path,
	}

	switch v := value.(type) {
	case map[string]any:
		for k, val := range v {
			childPath := k
			if path != "" {
				childPath = path + "." + k
			}
			child := buildTree(k, childPath, val)
			child.Parent = node
			node.Subsidiary = append(node.Subsidiary, child)
		}
	case []any:
		for i, val := range v {
			childKey := fmt.Sprintf("[%d]", i)
			childPath := path + childKey
			child := buildTree(childKey, childPath, val)
			node.Subsidiary = append(node.Subsidiary, child)
		}
	default:
		node.Value = v
	}

	return node
}
