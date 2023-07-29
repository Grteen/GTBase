package replic

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
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

		if !utils.EqualByteSlice(result, []byte(constants.HeartCommand+constants.CommandSep)) {
			t.Errorf("should read %v but got %v", constants.HeartCommand, result)
		}

		c := client.CreateGtBaseClient(fd, client.CreateAddress("127.0.0.1", port))
		m := replic.CreateMaster(10, 5000, c)

		errg := m.GetHeartFromMaster()
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

	result := make([]byte, 0)
	result = append(result, []byte(constants.GetHeartCommand)...)
	result = append(result, []byte(" ")...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(10)...)
	result = append(result, []byte(" ")...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(5000)...)
	// result = append(result, []byte(constants.CommandSep)...)
	if !utils.EqualByteSlice(res, result) {
		t.Errorf("should get %v but got %v", result, res)
	}
}
