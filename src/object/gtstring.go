package object

// GtString is used as Key and Value
type GtString struct {
	origin []byte
}

func (gts *GtString) Length() int {
	return len(gts.origin)
}
