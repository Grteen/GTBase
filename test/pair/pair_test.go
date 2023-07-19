package pair

import (
	"testing"
)

func TestPairToByte(t *testing.T) {
	data := []struct {
		key            string
		val            string
		flag           int8
		overFlowIndex  int32
		overFlowOffset int32
		res            string
	}{
		{"Key", "Val", 0, 0, 30, ""},
	}

	data = append(data)
}
