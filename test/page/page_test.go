package page

import (
	"GtBase/pkg/constants"
	"GtBase/src/nextwrite"
	"GtBase/src/object"
	"GtBase/src/page"
	"GtBase/src/pair"
	"os"
	"testing"

	"GtBase/utils"
)

func TestInitPageFile(t *testing.T) {
	page.InitPageFile()

	if _, err := os.Stat(constants.PageFilePathToDo); os.IsNotExist(err) {
		t.Errorf("InitPageFile() should create the %s but it didn't", constants.PageFilePathToDo)
	}
}

func TestReadWritePage(t *testing.T) {
	testReadWritePageInSingleIndex(t, 0)
	testReadWritePageInSingleIndex(t, 1)
	testReadWritePageInSingleIndex(t, 2)
}

func readWritePageCreateData() [][]byte {
	result := make([][]byte, 0)
	data := []string{"Hello World", "abc"}

	for _, d := range data {
		t := make([]byte, constants.PageSize)
		for i := 0; i < len(d); i++ {
			t[i] = d[i]
		}
		result = append(result, t)
	}

	return result
}

func testReadWritePageInSingleIndex(t *testing.T, idx int) {
	ph := page.CreatePageHeader(int32(idx))
	pg, err := page.ReadPage(ph.PageIndex())
	if err != nil {
		t.Errorf(err.Error())
	}

	data := readWritePageCreateData()

	for _, d := range data {
		pg.SetSrc(d)
		pg.FlushPage()

		spg, err := page.ReadPage(ph.PageIndex())
		if err != nil {
			t.Errorf(err.Error())
		}

		if !utils.EqualByteSlice(spg.Src(), d) {
			t.Errorf("WritePage should write %s but ReadPage reads %s", d, spg.Src())
		}
	}
}

func TestWriteBytes(t *testing.T) {
	page.DeletePageFile()
	page.InitPageFile()
	data := []struct {
		write []byte
		res   []byte
	}{
		{[]byte(""), []byte("")},
		{[]byte("First Write "), []byte("First Write ")},
		{[]byte("Second Write "), []byte("First Write Second Write ")},
		{[]byte("Hello World"), []byte("First Write Second Write Hello World")},
	}

	for i := 1; i < len(data); i++ {
		pg, err := page.ReadPage(1)
		if err != nil {
			t.Errorf(err.Error())
		}

		p, ok := page.GetPagePool().GetPage(1)
		if !ok {
			t.Errorf("GetPagePool should get index %v but not", 1)
		}

		if pg != p {
			t.Errorf("GetPagePool().GetPage() should be same as page.ReadBucketPage but not")
		}

		pg.WriteBytes(int32(len(data[i-1].res)), data[i].write)
		if pg.Dirty() != true {
			t.Errorf("page should be dirtied by WriteBytes but not")
		}

		if !utils.EqualByteSliceOnlyInMinLen(data[i].res, pg.Src()) {
			t.Errorf("page should be %v but it got %v", data[i].res, pg.Src()[:len(data[i].res)])
		}
	}

	pg, err := page.ReadPage(0)
	if err != nil {
		t.Errorf(err.Error())
	}

	pg.FlushPage()
	if pg.Dirty() != false {
		t.Errorf("page should be cleaned by FlushPage but not")
	}
}

func TestBucketWriteBytes(t *testing.T) {
	page.DeleteBucketPageFile()
	page.InitBucketPageFile()
	data := []struct {
		write []byte
		res   []byte
	}{
		{[]byte(""), []byte("")},
		{[]byte("First Write "), []byte("First Write ")},
		{[]byte("Second Write "), []byte("First Write Second Write ")},
		{[]byte("Hello World"), []byte("First Write Second Write Hello World")},
	}

	for i := 1; i < len(data); i++ {
		pg, err := page.ReadPage(-1)
		p, ok := page.GetPagePool().GetPage(-1)
		if !ok {
			t.Errorf("GetPagePool should get index %v but not", -1)
		}

		if pg != p {
			t.Errorf("GetPagePool().GetPage() should be same as page.ReadBucketPage but not")
		}

		if err != nil {
			t.Errorf(err.Error())
		}

		pg.WriteBytes(int32(len(data[i-1].res)), data[i].write)
		if pg.Dirty() != true {
			t.Errorf("page should be dirtied by WriteBytes but not")
		}

		if !utils.EqualByteSliceOnlyInMinLen(data[i].res, pg.Src()) {
			t.Errorf("page should be %v but it got %v", data[i].res, pg.Src()[:len(data[i].res)])
		}
	}

	pg, err := page.ReadPage(-1)
	if err != nil {
		t.Errorf(err.Error())
	}

	pg.FlushPage()
	if pg.Dirty() != false {
		t.Errorf("page should be cleaned by FlushPage but not")
	}
}

func TestReadPair(t *testing.T) {

	data := []struct {
		key         string
		val         string
		flag        int8
		overFlowIdx int32
		overFlowOff int32
	}{
		{"Key", "Hello", 1, 1, 5},
		{"Set", "Msg", 0, 15, 16383},
	}

	for _, d := range data {
		p := pair.CreatePair(object.CreateGtString(d.key), object.CreateGtString(d.val), d.flag, pair.CreateOverFlow(d.overFlowIdx, d.overFlowOff))
		nw, err := nextwrite.GetNextWriteAndIncreaseIt(int32(len(p.ToByte())))
		if err != nil {
			t.Errorf(err.Error())
		}

		idx, off := nw.NextWriteInfo()
		p.WriteInPage(idx, off)

		pg, errr := page.ReadPage(idx)
		if errr != nil {
			t.Errorf(err.Error())
		}

		midoff := int(off) + 4 + len(d.key) + 4 + len(d.val) + 1
		temp := midoff

		flag := pg.ReadFlag(int32(midoff))
		if flag != d.flag {
			t.Errorf("ReadFlag should got %v but got %v", d.flag, flag)
		}
		midoff -= 1

		keyLen := pg.ReadKeyLength(int32(midoff))
		if keyLen != int32(len(d.key)) {
			t.Errorf("ReadKeyLength should got %v but got %v", len(d.key), keyLen)
		}
		midoff -= 4

		key := pg.ReadKey(int32(midoff), keyLen)
		if key.ToString() != d.key {
			t.Errorf("ReadKey should got %v but got %v", d.key, key.ToString())
		}
		midoff -= int(keyLen)

		valLen := pg.ReadValLength(int32(midoff))
		if valLen != int32(len(d.val)) {
			t.Errorf("ReadValLength should got %v but got %v", len(d.val), valLen)
		}
		midoff -= 4

		val := pg.ReadVal(int32(midoff), valLen)
		if val.ToString() != d.val {
			t.Errorf("ReadVal should got %v but got %v", d.val, val)
		}

		midoff = temp

		overflowIdx, overflowOff := pg.ReadOverFlow(int32(midoff))
		if overflowIdx != d.overFlowIdx || overflowOff != d.overFlowOff {
			t.Errorf("ReadOverFlow should got %v idx %v off but got %v idx %v off", d.overFlowIdx, d.overFlowOff, overflowIdx, overflowOff)
		}
	}
}
