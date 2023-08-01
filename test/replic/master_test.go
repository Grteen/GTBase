package replic

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/page"
	"GtBase/src/replic"
	"GtBase/utils"
	"testing"
)

func TestMasterHeart(t *testing.T) {
	port := 6666
	lfd, err := utils.BindAndListen(port)
	if err != nil {
		t.Errorf(err.Error())
	}

	go func() {
		fd, err := utils.LocalDial(port)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer utils.CloseFd(fd)

		result, errr := utils.ReadFd(fd)
		if errr != nil {
			t.Errorf(errr.Error())
		}

		fields := make([][]byte, 0)
		fields = append(fields, []byte(constants.HeartCommand))
		fields = append(fields, utils.Encodeint32ToBytesSmallEnd(0))

		com := utils.EncodeFieldsToGtBasePacket(fields)

		if !utils.EqualByteSlice(result, com) {
			t.Errorf("should read %v but got %v", com, result)
		}

		c := client.CreateGtBaseClient(fd, client.CreateAddress("127.0.0.1", port))
		m := replic.CreateMaster(10, 5000, -1, c)

		errg := m.HeartFromMaster(0)
		if errg != nil {
			t.Errorf(errg.Error())
		}
	}()

	nfd, erra := utils.Accept(int(lfd))
	if erra != nil {
		t.Errorf(erra.Error())
	}

	c := client.CreateGtBaseClient(nfd, client.CreateAddress("127.0.0.1", port))
	s := replic.CreateSlave(0, 0, 1, c)

	errh := s.SendHeartToSlave()
	if errh != nil {
		t.Errorf(errh.Error())
	}

	res, errr := c.Read()
	if errr != nil {
		t.Errorf(errr.Error())
	}

	fields := make([][]byte, 0)
	fields = append(fields, []byte(constants.GetHeartCommand))
	fields = append(fields, utils.Encodeint32ToBytesSmallEnd(0))
	fields = append(fields, utils.Encodeint32ToBytesSmallEnd(10))
	fields = append(fields, utils.Encodeint32ToBytesSmallEnd(5000))

	result := utils.EncodeFieldsToGtBasePacket(fields)
	// result = append(result, []byte(constants.CommandSep)...)
	if !utils.EqualByteSlice(res, result[:len(result)-2]) {
		t.Errorf("should get %v but got %v", result[:len(result)-2], res)
	}
}

func TestRedo(t *testing.T) {
	port := 6544
	lfd, err := utils.BindAndListen(port)
	if err != nil {
		t.Errorf(err.Error())
	}

	redo, errr := page.ReadRedoPage(0)
	if errr != nil {
		t.Errorf(errr.Error())
	}

	slaveRedo := redo.Src()[:19]
	ch := make(chan struct{})
	go func() {
		fd, err := utils.LocalDial(port)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer utils.CloseFd(fd)

		result := make([]byte, 0)
		for {
			res, errr := utils.ReadFd(fd)
			if errr != nil {
				t.Errorf(errr.Error())
			}
			result = append(result, res...)
			if utils.EqualByteSlice(result[len(result)-2:], []byte(constants.ReplicRedoLogEnd)) {
				result = result[:len(result)-2]
				break
			}
		}
		seqBts := result[len(constants.RedoCommand)+1 : len(constants.RedoCommand)+1+int(constants.SendRedoLogSeqSize)]

		c := client.CreateGtBaseClient(fd, client.CreateAddress("127.0.0.1", port))
		m := replic.CreateMaster(0, int32(len(slaveRedo)), 1, c)

		errg := m.RedoFromMaster(utils.EncodeBytesSmallEndToint32(seqBts), result[len(constants.RedoCommand)+1+int(constants.SendRedoLogSeqSize)+1:])
		if errg != nil {
			t.Errorf(errg.Error())
		}
		ch <- struct{}{}
	}()

	nfd, erra := utils.Accept(int(lfd))
	if erra != nil {
		t.Errorf(erra.Error())
	}

	c := client.CreateGtBaseClient(nfd, client.CreateAddress("127.0.0.1", port))
	s := replic.CreateSlave(0, int32(len(slaveRedo)), 1, c)

	errh := s.SendRedoLogToSlave()
	if errh != nil {
		t.Errorf(errh.Error())
	}

	<-ch

	pg, err := page.ReadRedoPage(0)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !utils.EqualByteSlice(pg.Src(), redo.Src()) {
		t.Errorf("not same")
	}
}
