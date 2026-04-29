package parser

import (
	"configlinter/internal/domain"
	"fmt"

	"github.com/BurntSushi/toml"
)

type TOMLParser struct{}

func NewTOMLParser() *TOMLParser {
	return &TOMLParser{}
}

func (p *TOMLParser) Extensions() []string {
	return []string{".toml"}
}

func (p *TOMLParser) Parse(data []byte) (*domain.ConfigNode, error) {
	var raw any
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("toml parse error: %w", err)
	}
	root := buildTree("", "", raw)
	return root, nil
}
