package replic

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"GtBase/src/client"
	"GtBase/src/page"
	"GtBase/utils"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type NextSeq struct {
	seq        int32
	heartCount int32
	cLock      sync.Mutex
}

func (n *NextSeq) increaseCount() bool {
	n.cLock.Lock()
	defer n.cLock.Unlock()
	n.heartCount++

	if n.heartCount >= 3 {
		n.heartCount = 0
		return true
	}
	return false
}

func CreateNextSeq(seq int32) *NextSeq {
	return &NextSeq{seq: seq, heartCount: 0}
}

type HeartInfo struct {
	heartChan  chan int32
	heartCount int32
	heartSeq   int32
	hLock      sync.Mutex
}

// if heartCount greater than HeartCountLimit return true
func (h *HeartInfo) IncreaseCount() bool {
	h.hLock.Lock()
	defer h.hLock.Unlock()
	h.heartCount++
	if h.heartCount >= constants.HeartCountLimit {
		h.heartCount = 0
		return true
	}

	return false
}

func (h *HeartInfo) IncreaseSeq() {
	h.hLock.Lock()
	defer h.hLock.Unlock()
	h.heartSeq++
}

func (h *HeartInfo) Push(seq int32) {
	h.heartChan <- seq
}

func CreateHeartInfo() *HeartInfo {
	return &HeartInfo{heartChan: make(chan int32, 10)}
}

type Slave struct {
	client    *client.GtBaseClient
	logIdx    int32
	logOff    int32
	sLock     sync.Mutex
	nextSeq   *NextSeq
	syncState int32
	hf        *HeartInfo
}

func (s *Slave) GetSyncState() int32 {
	return s.syncState
}

func (s *Slave) SetSyncStateLock(state int32) {
	s.sLock.Lock()
	defer s.sLock.Unlock()
	s.syncState = state
}

func (s *Slave) GetClient() *client.GtBaseClient {
	return s.client
}

func (s *Slave) InitClient(host string, port int) error {
	fd, err := utils.Dial(host, port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fd, host, port)
	s.client = client.CreateGtBaseClient(fd, client.CreateAddress(host, port))
	// _, err = utils.WriteFd(fd, []byte("here\r\n"))
	return err
}

func (s *Slave) GetLogInfo() (int32, int32) {
	return s.logIdx, s.logOff
}

func (s *Slave) GetHeatInfo() *HeartInfo {
	return s.hf
}

// if get same Sequence three times it will return true
// and master should resend the redolog
func (s *Slave) GetSameSeq() bool {
	return s.nextSeq.increaseCount()
}

func (s *Slave) SetLogIdxAndOffLock(logIdx, logOff int32) {
	s.sLock.Lock()
	defer s.sLock.Unlock()
	s.logIdx = logIdx
	s.logOff = logOff
}

func (s *Slave) SetNextSeqLock(seq int32) {
	s.sLock.Lock()
	defer s.sLock.Unlock()
	s.nextSeq = CreateNextSeq(seq)
}

func (s *Slave) GetSeq() int32 {
	return s.nextSeq.seq
}

// return how many redo page can be send according to s.logIdx (not included s.logIdx page)
func (s *Slave) calRedoLogRestLen() (int32, error) {
	fileInfo, err := os.Stat(constants.RedoLogToDo)
	if err != nil {
		return -1, glog.Error("CalRedoLogRestLen can't Stat file %v becasuse %v", constants.RedoLogToDo, err)
	}

	fileSize := fileInfo.Size()
	totalLen := int32(fileSize) / int32(constants.PageSize)

	restPage := totalLen - s.logIdx - 1
	return restPage, nil
}

func (s *Slave) readRedoLogToSend(restPageLen int32) ([]byte, error) {
	if restPageLen == -1 {
		return make([]byte, 0), nil
	}
	firstPg, err := page.ReadRedoPage(s.logIdx)
	if err != nil {
		return nil, err
	}

	result := make([]byte, 0, (restPageLen+1)*int32(constants.PageSize))
	result = append(result, firstPg.SrcSlice(s.logOff, int32(constants.PageSize))...)

	for i := 1; i <= int(restPageLen); i++ {
		pg, err := page.ReadRedoPage(s.logIdx + int32(i))
		if err != nil {
			return nil, err
		}

		result = append(result, pg.Src()...)
	}

	return result, nil
}

// Redo seq redolog\r\n
func (s *Slave) SendRedoLogToSlave() error {
	restPageLen, err := s.calRedoLogRestLen()
	if err != nil {
		return err
	}

	pageToSendLen := restPageLen
	if pageToSendLen >= constants.MaxRedoLogPagesToSendOnce {
		pageToSendLen = constants.MaxRedoLogPagesToSendOnce
	}
	if pageToSendLen < -1 {
		pageToSendLen = -1
	}

	redoLog, errr := s.readRedoLogToSend(pageToSendLen)
	if errr != nil {
		return err
	}

	errw := client.Redo(s.client, redoLog, s.nextSeq.seq)
	if errw != nil {
		return errw
	}

	return nil
}

func (s *Slave) GetSendRedoLogResponseFromSlave(logIdx, logOff, seq int32) {
	s.SetLogIdxAndOffLock(logIdx, logOff)
	s.SetNextSeqLock(seq)
}

// return Slave's syncState and error
func (s *Slave) CheckFullSyncFinish() (int32, error) {
	restPageLen, err := s.calRedoLogRestLen()
	if err != nil {
		return -1, err
	}

	if restPageLen <= constants.SlaveFullSyncThreshold {
		s.SetSyncStateLock(constants.SlaveSync)
	}

	return s.syncState, nil
}

func (s *Slave) SendHeartToSlave() error {
	return client.Heart(s.client, s.hf.heartSeq)
}

func (s *Slave) GetHeartRespFromSlave(logIdx, logOff, seq, heartSeq int32) error {
	if heartSeq != s.hf.heartSeq {
		return nil
	}

	if seq <= s.GetSeq() {
		reSend := s.GetSameSeq()
		if reSend {
			if s.syncState == constants.SlaveFullSync {
				s.SetLogIdxAndOffLock(logIdx, logOff)
				s.SetNextSeqLock(seq)
				err := s.SendHeartToSlave()
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	if s.syncState == constants.SlaveSync {
		s.SetLogIdxAndOffLock(logIdx, logOff)
		s.SetNextSeqLock(seq)
		err := s.SendRedoLogToSlave()
		if err != nil {
			return err
		}
	}

	s.hf.Push(seq)

	return nil
}

func (s *Slave) HeartBeat(rs *ReplicState) {
loop:
	for {
		time.Sleep(1 * time.Second)
		s.SendHeartToSlave()
		select {
		case heartSeq := <-s.hf.heartChan:
			if heartSeq == s.hf.heartSeq {
				s.hf.IncreaseSeq()
			}
			goto loop
		case <-time.After(3 * time.Second):
			disc := s.hf.IncreaseCount()
			if disc {
				s.SetSyncStateLock(constants.SlaveDisConnect)
				return
			}
			goto loop
		}
	}
}

func CreateSlave(logIdx, logOff, seq int32, client *client.GtBaseClient) *Slave {
	return &Slave{logIdx: logIdx, logOff: logOff, nextSeq: CreateNextSeq(seq), syncState: constants.SlaveFullSync, hf: CreateHeartInfo(), client: client}
}
