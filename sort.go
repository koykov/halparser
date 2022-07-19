package halvector

// Custom implementation of quick sort algorithm, special for type []vector.Node (sub-slice of vector's node array).
// Need to avoid redundant allocation when using sort.Interface.
//
// sort.Interface problem:
// <code>
// type nodes []vector.Node // type that implements sort.Interface
// ...
// children := node.Children() // get a slice of nodes to sort
// nodes := nodes(children)    // <- simple typecast, but produces an alloc (copy) due to taking address in the next line
// sort.Sort(&nodes)           // taking address
// ...
// </code>

import "github.com/koykov/vector"

func pivot(vec *Vector, p []int, lo, hi int) int {
	if len(p) == 0 {
		return 0
	}
	pi := vec.NodeAt(p[hi])
	i := lo - 1
	_ = p[len(p)-1]
	for j := lo; j <= hi-1; j++ {
		pj := vec.NodeAt(p[j])
		if less(pj, pi) {
			i++
			vec.NodeAt(p[i]).SwapWith(pj)
		}
	}
	if i < hi {
		a, b := vec.NodeAt(p[i+1]), vec.NodeAt(p[hi])
		a.SwapWith(b)
	}
	return i + 1
}

func less(a, b *vector.Node) bool {
	aq, _ := a.DotFloat("quality")
	bq, _ := b.DotFloat("quality")
	return bq < aq
}

func quickSort(vec *Vector, p []int, lo, hi int) {
	if lo < hi {
		pi := pivot(vec, p, lo, hi)
		quickSort(vec, p, lo, pi-1)
		quickSort(vec, p, pi+1, hi)
	}
}

func (vec *Vector) Sort() *Vector {
	if vec.CheckBit(flagSorted) {
		return vec
	}
	vec.SetBit(flagSorted, true)
	ni := vec.Index.GetRow(1)
	quickSort(vec, ni, 0, len(ni)-1)
	return vec
}
