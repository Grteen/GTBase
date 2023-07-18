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

func (gts *GtString) ToByte() []byte {
	result := make([]byte, 0)
	copy(result, gts.origin)
	return result
}

func CreateGtString(str string) *GtString {
	gts := &GtString{}
	gts.SetValue(str)
	return gts
}
