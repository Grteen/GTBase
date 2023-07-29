package replic

import (
	"GtBase/src/client"
	"sync"
)

type Master struct {
	client *client.GtBaseClient
	logIdx int32
	logOff int32
	mLock  sync.Mutex
}

func (m *Master) GetClient() *client.GtBaseClient {
	return m.client
}

func (m *Master) GetLogInfo() (int32, int32) {
	return m.logIdx, m.logOff
}

func (m *Master) SetLogIdxAndOffLock(logIdx, logOff int32) {
	m.mLock.Lock()
	defer m.mLock.Unlock()
	m.logIdx = logIdx
	m.logOff = logOff
}

func (m *Master) sendGetHeartToMaster() error {
	return client.GetHeart(m.client, m.logIdx, m.logOff)
}

func (m *Master) HeartFromMaster() error {
	return m.sendGetHeartToMaster()
}

func CreateMaster(logIdx, logOff int32, client *client.GtBaseClient) *Master {
	return &Master{logIdx: logIdx, logOff: logOff, client: client}
}
