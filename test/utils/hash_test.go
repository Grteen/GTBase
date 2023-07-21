package utils

import (
	"GtBase/src/object"
	"GtBase/utils"
	"testing"
)

func TestHash(t *testing.T) {
	data := []string{"key", "val", "hello", "www.baidu.com-A", "www.tempfrost.kuocaitm.net"}
	for _, d := range data {
		f := utils.FirstHash(object.CreateGtString(d).ToByte())
		if f < 0 || f >= 256 {
			t.Errorf("FirstHash should got f > 0 and f < 256 but got %v", f)
		}

		s := utils.SecondHash(f)
		if s < 0 || s >= 256 {
			t.Errorf("FirstHash should got s > 0 and s < 256 but got %v", s)
		}
	}
}
