package domain

import "io/fs"

type ConfigNode struct {
	Key        string
	Path       string
	Value      any
	Subsidiary []*ConfigNode
	Parent     *ConfigNode
	FilePath   string
	FileMode   fs.FileMode
}

func (n *ConfigNode) IsScalar() bool {
	return len(n.Subsidiary) == 0
}

func (n *ConfigNode) StringValue() (string, bool) {
	s, ok := n.Value.(string)
	if ok {
		return s, true
	}
	return "", false
}

func (n *ConfigNode) BoolValue() (bool, bool) {
	if b, ok := n.Value.(bool); ok {
		return b, true
	}
	return false, false
}

func (n *ConfigNode) Walk(fn func(node *ConfigNode)) {
	fn(n)
	for _, child := range n.Subsidiary {
		child.Walk(fn)
	}
}
