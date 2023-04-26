package halvector

import "sync"

type Pool struct {
	p sync.Pool
}

var (
	P       Pool
	_, _, _ = Acquire, AcquireWithLimit, Release
)

func (p *Pool) Get() *Vector {
	v := p.p.Get()
	if v != nil {
		if vec, ok := v.(*Vector); ok {
			return vec
		}
	}
	return NewVector()
}

func (p *Pool) Put(vec *Vector) {
	vec.Reset()
	p.p.Put(vec)
}

func Acquire() *Vector {
	return P.Get()
}

func AcquireWithLimit(limit int) *Vector {
	return P.Get().SetLimit(limit)
}

func Release(vec *Vector) {
	P.Put(vec)
}
