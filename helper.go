package halvector

import (
	"io"

	"github.com/koykov/vector"
)

var helper Helper

type Helper struct{}

func (Helper) Indirect(p *vector.Byteptr) []byte          { return p.RawBytes() }
func (Helper) Beautify(_ io.Writer, _ *vector.Node) error { return nil }

func (Helper) Marshal(w io.Writer, node *vector.Node) error {
	_, err := w.Write(node.Bytes())
	return err
}
