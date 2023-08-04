package object

type GtNil struct {
}

func (obj *GtNil) ToByte() []byte {
	return []byte{0, 0, 0, 78, 105, 108}
}

func (obj *GtNil) ToString() string {
	return "Nil"
}

func CreateGtNil() Object {
	return &GtNil{}
}
