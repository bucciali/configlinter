package domain

import "io"

type Reporter interface {
	Report(w io.Writer, findings []Finding) error
}
