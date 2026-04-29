package domain

type Rule interface {
	ID() string
	Description() string
	Analyze(root *ConfigNode) []Finding
}
