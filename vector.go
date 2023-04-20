package halvector

import (
	"io"

	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

const (
	flagSorted = 8
)

type Vector struct {
	vector.Vector
	limit int
}

func NewVector() *Vector {
	vec := &Vector{}
	return vec
}

func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

func (vec *Vector) ParseStr(s string) error {
	return vec.parse(fastconv.S2B(s), false)
}

func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(fastconv.S2B(s), true)
}

// SetLimit setups hard of nodes. All entities over the limit will ignore.
// See BenchmarkLimit() for explanation.
func (vec *Vector) SetLimit(limit int) *Vector {
	if limit < 0 {
		limit = 0
	}
	vec.limit = limit
	return vec
}

func (vec *Vector) Beautify(w io.Writer) error {
	r := vec.Root()
	return vec.beautify(w, r, 0)
}

func (vec *Vector) Reset() {
	vec.Vector.Reset()
	vec.limit = 0
}
