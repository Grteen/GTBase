package replic

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"GtBase/src/client"
	"GtBase/src/page"
	"GtBase/utils"
	"os"
	"sync"
)

type NextSeq struct {
	seq   int32
	count int32
	cLock sync.Mutex
}

func (n *NextSeq) increaseCount() bool {
	n.cLock.Lock()
	defer n.cLock.Unlock()
	n.count++

	if n.count >= 3 {
		n.count = 0
		return true
	}
	return false
}

func CreateNextSeq(seq int32) *NextSeq {
	return &NextSeq{seq: seq, count: 0}
}

type Slave struct {
	client    *client.GtBaseClient
	logIdx    int32
	logOff    int32
	sLock     sync.Mutex
	nextSeq   *NextSeq
	syncState int32
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

func (s *Slave) GetLogInfo() (int32, int32) {
	return s.logIdx, s.logOff
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

// seq redolog \r\n
func (s *Slave) SendRedoLog() error {
	restPageLen, err := s.calRedoLogRestLen()
	if err != nil {
		return err
	}

	pageToSendLen := restPageLen
	if pageToSendLen >= constants.MaxRedoLogPagesToSendOnce {
		pageToSendLen = constants.MaxRedoLogPagesToSendOnce
	}

	result, errr := s.readRedoLogToSend(pageToSendLen)
	if errr != nil {
		return err
	}
	result = append(utils.Encodeint32ToBytesSmallEnd(s.nextSeq.seq), result...)
	result = append(result, []byte(constants.ReplicRedoLogEnd)...)

	errw := s.client.Write(result)
	if errw != nil {
		return errw
	}

	return nil
}

func (s *Slave) GetSendRedoLogResponse(logIdx, logOff, seq int32) {
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

func CreateSlave(logIdx, logOff, seq int32, client *client.GtBaseClient) *Slave {
	return &Slave{client: client, logIdx: logIdx, logOff: logOff, nextSeq: CreateNextSeq(seq), syncState: constants.SlaveFullSync}
}
