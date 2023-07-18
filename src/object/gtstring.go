package object

// GtString is used as Key and Value
type GtString struct {
	origin []byte
}

func (gts *GtString) Length() int {
	return len(gts.origin)
}

func (gts *GtString) SetValue(str string) {
	gts.origin = []byte(str)
}

func (gts *GtString) ToString() string {
	return string(gts.origin)
}
