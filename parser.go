package halvector

import (
	"errors"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

const (
	offsetCode    = 0
	offsetScript  = 4
	offsetRegion  = 10
	offsetQuality = 16
	offsetDefQT   = 23

	lenCode    = 4
	lenScript  = 6
	lenRegion  = 6
	lenQuality = 7
	lenDefQT   = 3
)

var (
	// Byte constants.
	bFmt   = []byte(" \t")
	bQt    = []byte(";q=")
	bComma = []byte(",")
	bSep   = []byte("-,")
	bKV    = []byte("codescriptregionquality1.0")

	ErrTooManyParts = errors.New("entry contains too many parts")
)

func (vec *Vector) parse(s []byte, copy bool) (err error) {
	s = bytealg.Trim(s, bFmt)
	if err = vec.SetSrc(s, copy); err != nil {
		return
	}

	offset := 0
	// Create root node and register it.
	root, i := vec.GetNodeWT(0, vector.TypeArr)
	root.SetOffset(vec.Index.Len(1))

	// Parse source data.
	offset, err = vec.parseGeneric(1, offset, root)
	if err != nil {
		vec.SetErrOffset(offset)
		return err
	}
	vec.PutNode(i, root)

	// Check unparsed tail.
	if offset < vec.SrcLen() {
		vec.SetErrOffset(offset)
		return vector.ErrUnparsedTail
	}

	return
}

func (vec *Vector) parseGeneric(depth, offset int, node *vector.Node) (int, error) {
	var (
		err error
		eof bool
	)
	for offset < vec.SrcLen() {
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		var qlo, qhi int
		if qlo = bytealg.IndexAt(vec.Src(), bQt, offset); qlo == -1 {
			qlo = vec.SrcLen()
		} else if qhi = bytealg.IndexAt(vec.Src(), bComma, qlo+3); qhi == -1 {
			qhi = vec.SrcLen()
		}
		if offset, err = vec.parseNode(depth, offset, qlo, qhi, node); err != nil {
			return offset, err
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, nil
		}
	}
	return offset, nil
}

func (vec *Vector) parseNode(depth, offset int, qlo, qhi int, root *vector.Node) (int, error) {
	if qhi < qlo {
		qhi = qlo
	}
	for {
		node, i := vec.GetChildWT(root, depth, vector.TypeObj)
		node.SetOffset(vec.Index.Len(depth + 1))
		p := bytealg.IndexAnyAt(vec.Src(), bSep, offset)
		if p == -1 {
			p = vec.SrcLen()
		}
		dc, d0, d1 := vec.indexDash(offset, qlo)
		if dc > 2 {
			return offset, ErrTooManyParts
		}

		switch dc {
		case 0:
			child, j := vec.GetChildWT(node, depth+1, vector.TypeStr)
			child.Key().Init(bKV, offsetCode, lenCode)
			child.Value().Init(vec.Src(), offset, qlo-offset)
			vec.PutNode(j, child)
			offset = qhi + 1
		case 1:
			child, j := vec.GetChildWT(node, depth+1, vector.TypeStr)
			child.Key().Init(bKV, offsetCode, lenCode)
			child.Value().Init(vec.Src(), offset, d0)
			vec.PutNode(j, child)
			offset = d0 + 1

			child, j = vec.GetChildWT(node, depth+1, vector.TypeStr)
			child.Key().Init(bKV, offsetRegion, lenRegion)
			child.Value().Init(vec.Src(), offset, qlo-offset)
			vec.PutNode(j, child)
			offset = qhi + 1
		case 2:
			// ...
		}

		// cl := qlo
		// if d0 != 0 {
		// 	cl = d0
		// }
		// child, j := vec.GetChildWT(node, depth+1, vector.TypeStr)
		// child.Key().Init(bKV, offsetCode, lenCode)
		// child.Value().Init(vec.Src(), offset, cl)
		// vec.PutNode(j, child)
		// offset = cl + 1

		_, _, _ = dc, d0, d1

		// if d1 > 0 {
		// 	child, j := vec.GetChildWT(node, depth+1, vector.TypeStr)
		// 	child.Key().Init(bKV, offsetScript, offsetScript+lenScript)
		// 	child.Value().Init(vec.Src(), offset, d1)
		// 	vec.PutNode(j, child)
		// 	offset = d1 + 1
		// }

		child, j := vec.GetChildWT(node, depth+1, vector.TypeNum)
		child.Key().Init(bKV, offsetQuality, lenQuality)
		if qlo > 0 && qhi > qlo {
			child.Value().Init(vec.Src(), qlo, qhi)
		} else {
			child.Value().Init(bKV, offsetDefQT, lenDefQT)
		}
		vec.PutNode(j, child)

		vec.PutNode(i, node)

		break
	}
	return offset, nil
}

func (vec *Vector) skipFmt(offset int) (int, bool) {
loop:
	if offset >= vec.SrcLen() {
		return offset, true
	}
	if c := vec.SrcAt(offset); c != bFmt[0] && c != bFmt[1] {
		return offset, false
	}
	offset++
	goto loop
}

func (vec *Vector) indexDash(lo, hi int) (dc, d0, d1 int) {
loop:
	if vec.SrcAt(lo) == '-' {
		dc++
		if dc == 1 {
			d0 = lo
		} else if dc == 2 {
			d1 = lo
		} else {
			return
		}
	}
	lo++
	if lo == hi {
		return
	}
	goto loop
}
