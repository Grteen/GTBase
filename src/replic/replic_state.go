package replic

import (
	"sync"
)

type ReplicState struct {
	slaves map[string]*Slave
	rsLock sync.Mutex
}

func (rs *ReplicState) AppendSlaveLock(s *Slave) {
	rs.rsLock.Lock()
	defer rs.rsLock.Unlock()

	key := s.client.GenerateKey()

	_, ok := rs.slaves[key]
	if !ok {
		rs.slaves[key] = s
	}
}

func (rs *ReplicState) GetSlave(key string) (*Slave, bool) {
	s, ok := rs.slaves[key]
	return s, ok
}

func CreateReplicState() *ReplicState {
	return &ReplicState{slaves: make(map[string]*Slave)}
}
