package nextset

// NextWrite tell the Next Write Command where to write
type NextWrite struct {
	pageIndex  int32
	pageOffset int32
}

func (nw *NextWrite) NextSetInfo() (int32, int32) {
	return nw.pageIndex, nw.pageOffset
}

type NextSetFactory struct {
}
