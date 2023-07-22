package object

// Object is the interface for the entire database storage corresponding to the data type used
type Object interface {
	ToString() string
	ToByte() []byte
}

func ParseObjectType(obj []byte) Object {
	return CreateGtString(string(obj))
}
