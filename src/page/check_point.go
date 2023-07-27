package page

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"GtBase/utils"
	"encoding/binary"
	"log"
	"os"
)

func InitCheckPointFile() {
	if _, err := os.Stat(constants.CheckPointPathToDo); os.IsNotExist(err) {
		file, errc := os.Create(constants.CheckPointPathToDo)
		if errc != nil {
			log.Fatalf("InitCheckPointFile can't create the CheckPointFile because %s\n", err)
		}

		errm := os.Chmod(constants.CheckPointPathToDo, 0777)
		if errm != nil {
			log.Fatalf("InitCheckPointFile can't chmod because of %s\n", errm)
		}

		errw := binary.Write(file, binary.LittleEndian, utils.Encodeint32ToBytesSmallEnd(0))
		if errw != nil {
			log.Fatalf("WriteCheckPointFile can't write file %v because %v", constants.CheckPointPathToDo, errw)
		}
	}
}

func DeleteCheckPointFile() {
	deletePageFile(constants.CheckPointPathToDo)
}

func WriteCheckPoint(cmn int32) error {
	if cmn <= 0 {
		return nil
	}

	return writeCheckPoint(cmn)
}

func writeCheckPoint(cmn int32) error {
	file, err := os.OpenFile(constants.CheckPointPathToDo, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return glog.Error("WriteCheckPoint can't open file %v because %v", constants.CheckPointPathToDo, err)
	}
	defer file.Close()

	errw := binary.Write(file, binary.LittleEndian, cmn)
	if errw != nil {
		return glog.Error("WriteCheckPoint can't write file %v because %v", constants.CheckPointPathToDo, errw)
	}

	return nil
}

func ReadCheckPoint() (int32, error) {
	file, err := os.OpenFile(constants.CheckPointPathToDo, os.O_RDWR, 0777)
	if err != nil {
		return -1, glog.Error("ReadCheckPoint can't open file %s because %s", constants.CheckPointPathToDo, err.Error())
	}
	defer file.Close()

	var result int32
	errr := binary.Read(file, binary.LittleEndian, &result)
	if errr != nil {
		return -1, glog.Error("ReadCheckPoint can't read file %s because %s", constants.CheckPointPathToDo, errr.Error())
	}

	return result, nil
}
