package halvector

import (
	"io"

	"github.com/koykov/vector"
)

var helper Helper

type Helper struct{}

func (Helper) Indirect(p *vector.Byteptr) []byte { return p.RawBytes() }

func (Helper) Beautify(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, true)
}

func (Helper) Marshal(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, false)
}
