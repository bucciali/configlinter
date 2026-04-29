package engine_test

import (
	"configlinter/internal/domain"
	"configlinter/internal/engine"
	"configlinter/internal/rules"
	"fmt"
	"testing"
)

func buildLargeTree(depth, breadth int) *domain.ConfigNode {
	root := &domain.ConfigNode{Key: "root", Path: "root"}
	buildChildren(root, depth, breadth)
	return root
}

func buildChildren(parent *domain.ConfigNode, depth, breadth int) {
	if depth == 0 {
		parent.Value = "some_value"
		return
	}
	for i := 0; i < breadth; i++ {
		child := &domain.ConfigNode{
			Key:    fmt.Sprintf("key_%d", i),
			Path:   parent.Path + fmt.Sprintf(".key_%d", i),
			Parent: parent,
		}
		parent.Subsidiary = append(parent.Subsidiary, child)
		buildChildren(child, depth-1, breadth)
	}
}

func BenchmarkEngine_SmallConfig(b *testing.B) {
	root := buildLargeTree(3, 3)
	e := engine.New(
		rules.NewPlaintextPasswordRule(),
		rules.NewBindAllRule(),
		rules.NewTLSDisabledRule(),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Analyze(root)
	}
}

func BenchmarkEngine_LargeConfig(b *testing.B) {
	root := buildLargeTree(5, 5)
	e := engine.New(
		rules.NewPlaintextPasswordRule(),
		rules.NewBindAllRule(),
		rules.NewTLSDisabledRule(),
		rules.NewDebugLogRule(),
		rules.NewWeakCryptoRule(),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Analyze(root)
	}
}
