package replic

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/page"
	"GtBase/utils"
	"sync"
)

type Master struct {
	client *client.GtBaseClient
	logIdx int32
	logOff int32
	seq    int32
	mLock  sync.Mutex
}

func (m *Master) GetClient() *client.GtBaseClient {
	return m.client
}

func (m *Master) GetLogInfo() (int32, int32) {
	return m.logIdx, m.logOff
}

func (m *Master) GetSeq() int32 {
	return m.seq
}

func (m *Master) SetLogIdxAndOffLock(logIdx, logOff int32) {
	m.mLock.Lock()
	defer m.mLock.Unlock()
	m.logIdx = logIdx
	m.logOff = logOff
}

func (m *Master) SetSeqLock(seq int32) {
	m.mLock.Lock()
	defer m.mLock.Unlock()
	m.seq = seq
}

func (m *Master) sendGetHeartToMaster(heartSeq int32) error {
	return client.GetHeart(m.client, m.logIdx, m.logOff, heartSeq)
}

func (m *Master) HeartFromMaster(heartSeq int32) error {
	return m.sendGetHeartToMaster(heartSeq)
}

func (m *Master) sendGetRedoToMaster() error {
	return client.GetRedo(m.client, m.logIdx, m.logOff, m.seq)
}

func (m *Master) RedoFromMaster(seq int32, redoLog []byte) (*utils.Message, error) {
	if seq != m.seq {
		return nil, nil
	}

	err := page.WriteRedoLogFromReplic(m.logIdx, m.logOff, redoLog)
	if err != nil {
		return nil, err
	}
	m.updateLogIdxAndOff(int32(len(redoLog)))
	m.SetSeqLock(seq + 1)
	return utils.CreateMessage(constants.MessageNeedRedo), m.sendGetRedoToMaster()
}

func (m *Master) updateLogIdxAndOff(redoLogLen int32) {
	restLen := redoLogLen
	tempIdx, tempOff := m.GetLogInfo()
	tempOff += restLen
	tempIdx += tempOff / int32(constants.PageSize)
	tempOff %= int32(constants.PageSize)

	m.SetLogIdxAndOffLock(tempIdx, tempOff)
}

func CreateMaster(logIdx, logOff, seq int32, client *client.GtBaseClient) *Master {
	return &Master{logIdx: logIdx, logOff: logOff, client: client, seq: seq}
}
