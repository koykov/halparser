package halvector

func skipFmt(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	if src[offset] > ' ' {
		return offset, false
	}
	for ; offset < n; offset++ {
		c := src[offset]
		if c != ' ' && c != '\t' {
			return offset, false
		}
	}
	return offset, true
}

func skipFmtTable(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	if offset == n {
		return offset, true
	}
	if src[offset] > ' ' {
		return offset, false
	}
	_ = skipTable[255]
	for ; skipTable[src[offset]]; offset++ {
	}
	return offset, offset == n
}

var skipTable = [256]bool{}

func init() {
	skipTable[' '] = true
	skipTable['\t'] = true
}

var _ = skipFmt
