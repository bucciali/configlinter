package domain

type Parser interface {
	Extensions() []string
	Parse(data []byte) (*ConfigNode, error)
}
