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
	bQt = []byte(";q=")
	bKV = []byte("codescriptregionquality1.0")

	ErrTooManyParts = errors.New("entry contains too many parts")
)

func (vec *Vector) parse(s []byte, copy bool) (err error) {
	if vec.Helper == nil {
		vec.Helper = helper
	}

	s = bytealg.TrimBytesFmt4(s)
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
		c   int
	)
	src := vec.Src()[offset:]
	n := len(src)
	_ = src[n-1]
	for offset < n {
		if offset, eof = skipFmtTable(src, n, offset); eof {
			return offset, vector.ErrUnexpEOF
		}

		var nhi int
		if nhi = bytealg.IndexByteAtBytes(src, ',', offset); nhi == -1 {
			nhi = n
		}

		var qlo, qhi int
		if qlo = bytealg.IndexAtBytes(src[:nhi], bQt, offset); qlo == -1 {
			qlo = nhi
		} else {
			qhi = nhi
		}
		if offset, err = vec.parseNode(depth, offset, qlo, qhi, node); err != nil {
			return offset, err
		}
		c++
		if vec.limit > 0 && c >= vec.limit {
			// Replace offset to SrcLen to avoid unparsed tail error.
			return n, nil
		}
		if offset, eof = skipFmtTable(src, n, offset); eof {
			return offset, nil
		}
		if src[offset] == ',' {
			if offset+1 < n && src[offset+1] == ';' {
				// Detect broken format, see testdata/15.hal.txt for example.
				return offset, vector.ErrUnexpId
			}
			offset++
		}
		if offset, eof = skipFmtTable(src, n, offset); eof {
			return offset, nil
		}
	}
	return offset, nil
}

func (vec *Vector) parseNode(depth, offset int, qlo, qhi int, root *vector.Node) (int, error) {
	var eof bool
	if qhi < qlo {
		qhi = qlo
	}
	src := vec.Src()
	n := len(src)
	_ = src[n-1]
	for {
		if offset == qlo {
			break
		}

		node, i := vec.GetChildWT(root, depth, vector.TypeObj)
		node.SetOffset(vec.Index.Len(depth + 1))
		p := bytealg.IndexByteAtBytes(src, '-', offset)
		if p == -1 {
			p = n
		}
		dc, d0, d1 := indexDash(src, n, offset, qlo)
		if dc > 2 {
			return offset, ErrTooManyParts
		}

		switch dc {
		case 0:
			// Add only code.
			vec.addStrNode(node, depth+1, offsetCode, lenCode, offset, qlo-offset)
			offset = qlo
		case 1:
			// Add code and region.
			vec.addStrNode(node, depth+1, offsetCode, lenCode, offset, d0-offset)
			offset = d0 + 1
			vec.addStrNode(node, depth+1, offsetRegion, lenRegion, offset, qlo-offset)
			offset = qlo
		case 2:
			// Add code, script and region.
			vec.addStrNode(node, depth+1, offsetCode, lenCode, offset, d0-offset)
			offset = d0 + 1
			vec.addStrNode(node, depth+1, offsetScript, lenScript, offset, d1-offset)
			offset = d1 + 1
			vec.addStrNode(node, depth+1, offsetRegion, lenRegion, offset, qlo-offset)
			offset = qlo
		}

		// Write quality.
		child, j := vec.GetChildWT(node, depth+1, vector.TypeNum)
		child.Key().Init(bKV, offsetQuality, lenQuality)
		if qlo > 0 && qhi > qlo {
			child.Value().Init(src, qlo+3, qhi-(qlo+3)) // +3 means length of ";q="
		} else {
			child.Value().Init(bKV, offsetDefQT, lenDefQT)
		}
		vec.PutNode(j, child)

		vec.PutNode(i, node)

		if offset, eof = skipFmtTable(src, n, offset); eof {
			return offset, nil
		}
		if offset == qlo {
			offset = qhi
			break
		}
		if offset, eof = skipFmtTable(src, n, offset); eof {
			return offset, nil
		}
	}
	return offset, nil
}

func (vec *Vector) addStrNode(root *vector.Node, depth, kpos, klen, vpos, vlen int) {
	node, j := vec.GetChildWT(root, depth, vector.TypeStr)
	node.Key().Init(bKV, kpos, klen)
	node.Value().Init(vec.Src(), vpos, vlen)
	vec.PutNode(j, node)
}

func indexDash(src []byte, n, lo, hi int) (dc, d0, d1 int) {
	_ = src[n-1]
loop:
	if src[lo] == '-' {
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
