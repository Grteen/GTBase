package replic

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/nextwrite"
	"GtBase/src/page"
	"GtBase/utils"
	"log"
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
	s.client = client.CreateGtBaseClient(fd, client.CreateAddress(host, port))
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

func (s *Slave) calRedoLogRestLen() (int64, error) {
	nw, err := nextwrite.GetRedoNextWriteAndIncreaseIt(0)
	if err != nil {
		return -1, err
	}

	nwidx, nwoff := nw.NextWriteInfo()

	result := (int64(nwidx)*constants.PageSize + int64(nwoff)) - (int64(s.logIdx)*constants.PageSize + int64(s.logOff))
	if result > constants.MaxRedoLogToSendOnceint32 {
		return constants.MaxRedoLogToSendOnceint32, nil
	}

	if result < 0 {
		result = 0
	}

	return result, nil
}

func (s *Slave) readRedoLogToSend() ([]byte, error) {
	restLen, errc := s.calRedoLogRestLen()
	if errc != nil {
		return nil, errc
	}

	firstPg, err := page.ReadRedoPage(s.logIdx)
	if err != nil {
		return nil, err
	}
	if restLen+int64(s.logOff) < constants.PageSize {
		return firstPg.SrcSliceLength(s.logOff, int32(restLen)), nil
	}

	rest := restLen
	result := make([]byte, 0)
	result = append(result, firstPg.SrcSlice(s.logOff, int32(constants.PageSize))...)
	rest -= (constants.PageSize - int64(s.logOff))
	tempIdx := s.logIdx
	tempIdx += 1
	for {
		pg, err := page.ReadRedoPage(tempIdx)
		if err != nil {
			return nil, err
		}

		if rest >= constants.PageSize {
			result = append(result, pg.Src()...)
			tempIdx += 1
			rest -= constants.PageSize
		} else {
			result = append(result, pg.SrcSliceLength(0, int32(restLen))...)
			break
		}
	}

	return result, nil
}

// Redo seq redolog\r\n
func (s *Slave) SendRedoLogToSlave(uuid string) error {
	redoLog, err := s.readRedoLogToSend()
	if err != nil {
		return err
	}

	errw := client.Redo(s.client, redoLog, s.nextSeq.seq, uuid)
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
	restLen, err := s.calRedoLogRestLen()
	if err != nil {
		return -1, err
	}

	if restLen <= constants.SlaveFullSyncThreshold {
		s.SetSyncStateLock(constants.SlaveSync)
	}

	return s.syncState, nil
}

func (s *Slave) SendHeartToSlave(uuid string) error {
	return client.Heart(s.client, s.hf.heartSeq, uuid)
}

func (s *Slave) GetHeartRespFromSlave(logIdx, logOff, seq, heartSeq int32, uuid string, uuidSelf string) error {
	if heartSeq != s.hf.heartSeq {
		return nil
	}

	if seq <= s.GetSeq() {
		reSend := s.GetSameSeq()
		if reSend {
			if s.syncState == constants.SlaveFullSync {
				s.SetLogIdxAndOffLock(logIdx, logOff)
				s.SetNextSeqLock(seq)
				err := s.SendHeartToSlave(uuidSelf)
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
		err := s.SendRedoLogToSlave(uuid)
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
		s.SendHeartToSlave(rs.GetUUID())
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
