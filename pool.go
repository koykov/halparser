package halvector

import (
	"sync"

	"github.com/koykov/vector"
)

type Pool struct {
	p sync.Pool
}

var p Pool

// Get old vector from the pool or create new one.
func (p *Pool) Get() *Vector {
	v := p.p.Get()
	if v != nil {
		if vec, ok := v.(*Vector); ok {
			return vec
		}
	}
	return NewVector()
}

// Put vector back to the pool.
func (p *Pool) Put(vec *Vector) {
	vec.Reset()
	p.p.Put(vec)
}

// Acquire returns vector from default pool instance.
func Acquire() *Vector {
	return p.Get()
}

// AcquireWithLimit returns vector with given limit.
func AcquireWithLimit(limit int) *Vector {
	return p.Get().SetLimit(limit)
}

// Release puts vector back to default pool instance.
func Release(vec *Vector) {
	p.Put(vec)
}

// ReleaseNC puts vector back to pool with enforced no-clear flag.
func ReleaseNC(vec *Vector) {
	vec.SetBit(vector.FlagNoClear, true)
	p.Put(vec)
}

var _, _, _, _ = Acquire, AcquireWithLimit, Release, ReleaseNC
