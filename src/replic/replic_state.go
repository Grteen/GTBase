package replic

import (
	"sync"
)

type ReplicState struct {
	slaves map[string]*Slave
	sLock  sync.Mutex
	master *Master
	mLock  sync.Mutex
}

func (rs *ReplicState) AppendSlaveLock(s *Slave) {
	rs.sLock.Lock()
	defer rs.sLock.Unlock()

	key := s.client.GenerateKey()

	_, ok := rs.slaves[key]
	if !ok {
		rs.slaves[key] = s
	}
}

func (rs *ReplicState) SetMasterLock(m *Master) {
	rs.mLock.Lock()
	defer rs.mLock.Unlock()

	rs.master = m
}

func (rs *ReplicState) GetSlave(key string) (*Slave, bool) {
	s, ok := rs.slaves[key]
	return s, ok
}

func (rs *ReplicState) GetMaster() *Master {
	return rs.master
}

func CreateReplicState() *ReplicState {
	return &ReplicState{slaves: make(map[string]*Slave)}
}
