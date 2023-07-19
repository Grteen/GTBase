package nextset

// NextSet tell the Next Set Command where to set
type NextSet struct {
	pageIndex  int32
	pageOffset int32
}

func (ns *NextSet) NextSetInfo() (int32, int32) {
	return ns.pageIndex, ns.pageOffset
}
